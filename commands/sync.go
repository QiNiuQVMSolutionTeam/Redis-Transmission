package commands

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Synchronizer struct {
	Workers map[uint64]*SyncWorker
}

type SyncWorker struct {
	DatabaseId              uint64
	SourceClient            *redis.Client
	DestinationClient       *redis.Client
	KeysPipeline            chan string
	DestinationKeysPipeline chan string
}

func (s *Synchronizer) InitClients(sourceHost, sourcePassword, destinationHost, destinationPassword string, dbCount uint64) {

	s.Workers = make(map[uint64]*SyncWorker, dbCount)

	for dbId := uint64(0); dbId < dbCount; dbId++ {

		s.Workers[dbId] = &SyncWorker{
			DatabaseId: dbId,
			SourceClient: redis.NewClient(&redis.Options{
				Addr:     sourceHost,
				Password: sourcePassword, // no password set
				DB:       int(dbId),      // use default DB
			}),
			DestinationClient: redis.NewClient(&redis.Options{
				Addr:     destinationHost,
				Password: destinationPassword, // no password set
				DB:       int(dbId),           // use default DB
			}),
		}
	}
}

func (s *Synchronizer) Go() {

	log.Println("Starting synchronizer")
	for _, worker := range s.Workers {

		go func(worker *SyncWorker) {
			for {
				if worker.Sync() <= 0 {

					time.Sleep(time.Second)
				}
			}
		}(worker)
	}

	for {
		time.Sleep(time.Hour)
	}
}

func (w *SyncWorker) Sync() (count uint64) {

	log.Printf("Start %d database thread\n", w.DatabaseId)
	w.InitChannel()
	go w.ReadKeys()

	count = w.WriteData()

	go w.ReadDestinationKeys()
	count += w.CheckNotExistKeys()

	log.Printf("Synchronized database(%d) %d records.", w.DatabaseId, count)

	return
}

func (w *SyncWorker) InitChannel() {

	w.KeysPipeline = make(chan string)
	w.DestinationKeysPipeline = make(chan string)
}

func (w *SyncWorker) ReadKeys() {

	var currentCursor uint64
	for {

		keys, nextCursor, err := w.SourceClient.Scan(currentCursor, "", 100).Result()

		if err != nil {

			log.Printf("Scan database(%d) error , %s\n", currentCursor, err)
			break
		}

		for _, key := range keys {

			w.KeysPipeline <- key
		}

		if nextCursor == 0 {

			break
		}

		currentCursor = nextCursor
	}

	close(w.KeysPipeline)
}

func (w *SyncWorker) WriteData() (count uint64) {

	for {
		key := w.getKey()
		if key == "" {
			break
		}

		record, err := w.dump(key)
		if err != nil {
			log.Printf("Dump key \"%s\" error, %s\n", key, err)
			continue
		}
		err = w.restore(record)
		if err != nil {
			log.Printf("Restore key \"%s\" error, %s\n", key, err)
			continue
		}

		count++
	}

	return
}

func (w *SyncWorker) getKey() string {

	var key string

	for {
		select {
		case key = <-w.KeysPipeline:

			return key

		default:

			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
}

func (w *SyncWorker) dump(key string) (record TransferRecord, err error) {

	record.Key = key
	record.TTL, err = w.SourceClient.TTL(key).Result()
	if err != nil {

		return
	}

	record.Value, err = w.SourceClient.Dump(key).Result()
	if err != nil {

		return
	}

	return
}

func (w *SyncWorker) restore(record TransferRecord) (err error) {

	duration, err := time.ParseDuration(fmt.Sprintf("%ds", record.TTL))
	if err != nil {
		log.Printf("Parse ttl(%d) error, %s\n", record.TTL, err)
		return
	}

	if duration > 0 {
		_, err = w.DestinationClient.RestoreReplace(record.Key, duration, record.Value).Result()
	} else {
		_, err = w.DestinationClient.RestoreReplace(record.Key, 0, record.Value).Result()
	}

	return
}

func (w *SyncWorker) ReadDestinationKeys() {

	var currentCursor uint64
	for {

		keys, nextCursor, err := w.DestinationClient.Scan(currentCursor, "", 100).Result()

		if err != nil {

			log.Printf("Scan destination database(%d) error , %s\n", currentCursor, err)
			break
		}

		for _, key := range keys {

			w.DestinationKeysPipeline <- key
		}

		if nextCursor == 0 {

			break
		}

		currentCursor = nextCursor
	}

	close(w.DestinationKeysPipeline)
}

func (w *SyncWorker) CheckNotExistKeys() (count uint64) {

	for {
		key := w.getDestinationKey()
		if key == "" {
			break
		}

		if w.sourceExist(key) {

			continue
		}

		err := w.removeDestinationKey(key)
		if err != nil {
			log.Printf("Remove key \"%s\" error, %s\n", key, err)
			continue
		}

		count++
	}

	return
}

func (w *SyncWorker) getDestinationKey() string {

	var key string

	for {
		select {
		case key = <-w.DestinationKeysPipeline:

			return key

		default:

			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
}

func (w *SyncWorker) sourceExist(key string) bool {

	isExist, err := w.SourceClient.Exists(key).Result()
	if err != nil {
		log.Printf("Judge Key in source error , key: %s , error: %s\n", key, err)
		return true
	}

	return isExist != 0
}

func (w *SyncWorker) removeDestinationKey(key string) (err error) {

	_, err = w.DestinationClient.Del(key).Result()
	return
}

func Sync(sourceHost, sourcePassword, destinationHost, destinationPassword string, databaseCount uint64) {

	s := &Synchronizer{}
	if databaseCount == 0 {
		databaseCount = getDatabaseCount(sourceHost, sourcePassword)
	}

	if databaseCount == 0 {

		log.Println("Get database count error.")
		return
	}

	s.InitClients(sourceHost, sourcePassword, destinationHost, destinationPassword, databaseCount)
	s.Go()
}
