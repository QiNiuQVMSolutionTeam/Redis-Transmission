package commands

import (
	"github.com/QiNiuQVMSolutionTeam/Redis-Transmission/lib"
	"github.com/go-redis/redis"
	"go.uber.org/atomic"
	"log"
	"sync"
	"time"
)

type Synchronizer struct {
	Workers map[uint64]*SyncOneRound
}

type SyncOneRound struct {
	DatabaseId              uint64
	SourceClient            *redis.Client
	DestinationClient       *redis.Client
	KeysPipeline            chan string
	DestinationKeysPipeline chan string
	Workers                 *lib.Workers
	ThreadCount             int
}

type SyncWorker struct {
	SourceClient      *redis.Client
	DestinationClient *redis.Client
}

func (s *Synchronizer) InitClients(sourceHost, sourcePassword, destinationHost, destinationPassword string, dbCount uint64, threadCount int) {

	s.Workers = make(map[uint64]*SyncOneRound, dbCount)

	for dbId := uint64(0); dbId < dbCount; dbId++ {

		s.Workers[dbId] = &SyncOneRound{
			DatabaseId: dbId,
			SourceClient: redis.NewClient(&redis.Options{
				Addr:     sourceHost,
				Password: sourcePassword,
				DB:       int(dbId),
				PoolSize: threadCount,
			}),
			DestinationClient: redis.NewClient(&redis.Options{
				Addr:     destinationHost,
				Password: destinationPassword,
				DB:       int(dbId),
				PoolSize: threadCount,
			}),
			ThreadCount: threadCount,
		}
	}
}

func (s *Synchronizer) Go(syncTimes uint64) {

	var wg sync.WaitGroup
	log.Println("Starting synchronizer")
	for _, worker := range s.Workers {

		wg.Add(1)
		go func(worker *SyncOneRound, syncTimes uint64) {
			for {
				if worker.Sync() <= 0 {

					time.Sleep(time.Second)
				}

				if syncTimes > 0 {

					syncTimes--
					if syncTimes <= 0 {
						break
					}
				}
			}

			wg.Done()
		}(worker, syncTimes)
	}

	wg.Wait()
}

func (round *SyncOneRound) Sync() (count uint64) {

	log.Printf("Start %d database thread\n", round.DatabaseId)
	round.InitChannel()
	go round.ReadKeys()

	count = round.SyncData()

	go round.ReadDestinationKeys()
	count += round.CheckNotExistKeys()

	log.Printf("Synchronized database(%d) %d records.", round.DatabaseId, count)

	return
}

func (round *SyncOneRound) InitChannel() {

	round.KeysPipeline = make(chan string, 1000)
	round.DestinationKeysPipeline = make(chan string, 1000)
	round.Workers = lib.NewWorkers(round.ThreadCount, func() interface{} {
		return &SyncWorker{
			SourceClient:      round.SourceClient,
			DestinationClient: round.DestinationClient,
		}
	})
}

func (round *SyncOneRound) ReadKeys() {

	log.Printf("Scan database(%d) start\n", round.DatabaseId)
	var currentCursor, keyCount uint64
	for {

		keys, nextCursor, err := round.SourceClient.Scan(currentCursor, "", 1000).Result()

		if err != nil {

			log.Printf("Scan database(%d) error , %s\n", currentCursor, err)
			break
		}

		for _, key := range keys {

			round.KeysPipeline <- key
		}

		if nextCursor == 0 {

			break
		}

		currentCursor = nextCursor
		keyCount += uint64(len(keys))
	}

	close(round.KeysPipeline)
	log.Printf("Scan database(%d) finished\n", round.DatabaseId)
}

func (round *SyncOneRound) SyncData() (uint64) {

	var count atomic.Uint64
	for {
		key := round.getKey()
		if key == "" {
			break
		}

		worker := round.getWorker()

		go func(key string) {

			defer func() {
				round.putWorker(worker)
			}()

			record, err := worker.dump(key)
			if err != nil {
				log.Printf("Dump key \"%s\" error, %s\n", key, err)
				return
			}
			err = worker.restore(record)
			if err != nil {
				log.Printf("Restore key \"%s\" error, %s\n", key, err)
				return
			}

			count.Inc()
		}(key)
	}

	round.Workers.Wait()

	return count.Load()
}

func (round *SyncOneRound) getKey() string {

	var key string

	for {
		select {
		case key = <-round.KeysPipeline:

			return key

		default:

			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
}

func (round *SyncOneRound) ReadDestinationKeys() {

	log.Printf("Scan destination database(%d) start\n", round.DatabaseId)
	var currentCursor uint64
	for {

		keys, nextCursor, err := round.DestinationClient.Scan(currentCursor, "", 100).Result()

		if err != nil {

			log.Printf("Scan destination database(%d) error , %s\n", currentCursor, err)
			break
		}

		for _, key := range keys {

			round.DestinationKeysPipeline <- key
		}

		if nextCursor == 0 {

			break
		}

		currentCursor = nextCursor
	}

	close(round.DestinationKeysPipeline)
	log.Printf("Scan destination database(%d) finished\n", round.DatabaseId)
}

func (round *SyncOneRound) CheckNotExistKeys() (uint64) {

	var count atomic.Uint64
	for {
		key := round.getDestinationKey()
		if key == "" {
			break
		}

		worker := round.getWorker()

		go func(key string) {

			defer func() {
				round.putWorker(worker)
			}()

			if worker.sourceExist(key) {

				return
			}

			err := worker.removeDestinationKey(key)
			if err != nil {
				log.Printf("Remove key \"%s\" error, %s\n", key, err)
				return
			}

			count.Inc()
		}(key)
	}

	round.Workers.Wait()

	return count.Load()
}

func (round *SyncOneRound) getDestinationKey() string {

	var key string

	for {
		select {
		case key = <-round.DestinationKeysPipeline:

			return key

		default:

			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
}

func (round *SyncOneRound) getWorker() *SyncWorker {

	return round.Workers.Get().(*SyncWorker)
}

func (round *SyncOneRound) putWorker(worker *SyncWorker) {

	round.Workers.Put(worker)
}

func (round *SyncWorker) dump(key string) (record TransferRecord, err error) {

	record.Key = key
	record.TTL, err = round.SourceClient.TTL(key).Result()
	if err != nil {

		return
	}

	record.Value, err = round.SourceClient.Dump(key).Result()
	if err != nil {

		return
	}

	return
}

func (round *SyncWorker) restore(record TransferRecord) (err error) {

	if record.TTL > 0 {
		_, err = round.DestinationClient.RestoreReplace(record.Key, record.TTL, record.Value).Result()
	} else {
		_, err = round.DestinationClient.RestoreReplace(record.Key, 0, record.Value).Result()
	}

	return
}

func (round *SyncWorker) sourceExist(key string) bool {

	isExist, err := round.SourceClient.Exists(key).Result()
	if err != nil {
		log.Printf("Judge Key in source error , key: %s , error: %s\n", key, err)
		return true
	}

	return isExist != 0
}

func (round *SyncWorker) removeDestinationKey(key string) (err error) {

	_, err = round.DestinationClient.Del(key).Result()
	return
}

func Sync(sourceHost, sourcePassword, destinationHost, destinationPassword string, databaseCount, syncTimes uint64, threadCount int) {

	s := &Synchronizer{}
	if databaseCount == 0 {
		databaseCount = getDatabaseCount(sourceHost, sourcePassword)
	}

	if databaseCount == 0 {

		log.Println("Get database count error.")
		return
	}

	s.InitClients(sourceHost, sourcePassword, destinationHost, destinationPassword, databaseCount, threadCount)
	s.Go(syncTimes)
}
