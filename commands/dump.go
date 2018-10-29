package commands

import (
	"encoding/base64"
	"fmt"
	"github.com/QiNiuQVMSolutionTeam/Redis-Transmission/lib"
	"github.com/gin-gonic/gin/json"
	"github.com/go-redis/redis"
	"go.uber.org/atomic"
	"log"
	"os"
	"strconv"
)

type Dumper struct {
	Host        string
	Password    string
	Client      *redis.Client
	Path        string
	DatabaseId  uint64
	Count       atomic.Uint64
	workers     *lib.Workers
	hasError    bool
	Stream      *os.File
	ThreadCount int
}

type DumpWorker struct {
	Client     *redis.Client
	DatabaseId uint64
	stream     *os.File
}

func (d *Dumper) Dump() {

	d.initSemaphore(d.ThreadCount)

	cursor := uint64(0)

	for {

		keys, nextCursor, err := d.scan(cursor)
		if err != nil {

			log.Printf("Error: Scan keys error, %s\n", err)
			break
		}

		for _, key := range keys {

			worker := d.getSemaphore()

			go func(key string) {

				defer func() {
					d.putSemaphore(worker)
				}()
				err := worker.Dump(key)

				if err != nil {
					d.hasError = true
					return
				}

				d.Count.Inc()

				if d.Count.Load()%1000 == 0 {
					d.PrintReport()
				}

			}(key)

			if d.hasError {
				break
			}
		}

		d.workers.Wait()

		if d.hasError {
			break
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	d.CloseClient()
	d.closeSemaphore()

	d.PrintReport()
}

func (d *Dumper) scan(cursor uint64) (keys []string, nextCursor uint64, err error) {

	keys, nextCursor, err = d.Client.Scan(cursor, "", 100).Result()
	return
}

func (d *Dumper) CloseClient() {

	if _, err := d.Client.Ping().Result(); err != nil {
		return
	}

	d.Client.Close()
}

func (d *Dumper) PrintReport() {

	log.Printf("DB %d dumped %d Record(s).\n", d.DatabaseId, d.Count.Load())
}

func (d *Dumper) initSemaphore(threadCount int) {

	d.workers = lib.NewWorkers(threadCount,
		func() interface{} {
			return &DumpWorker{
				Client:     d.Client,
				DatabaseId: d.DatabaseId,
				stream:     d.Stream,
			}
		},
	)
}

func (d *Dumper) getSemaphore() *DumpWorker {

	dw := d.workers.Get()
	return dw.(*DumpWorker)
}

func (d *Dumper) putSemaphore(dw *DumpWorker) {

	d.workers.Put(dw)
}

func (d *Dumper) closeSemaphore() {

	d.workers.Wait()
	for d.workers.IdleCount() > 0 {

		worker := d.getSemaphore()
		if worker == nil {

			break
		}

		worker.CloseClient()
	}
}

func (dw *DumpWorker) Dump(key string) (err error) {

	record := &Record{Key: key}

	record.Value, err = dw.getSerializeString(key)

	if err != nil {

		log.Printf("Error: Get key serialize string error, %s\n", err)
		log.Printf("Key: %s\n", key)
		log.Printf("Client: %#v\n", dw.Client)
		log.Printf("Pool: %#v\n", *dw.Client.PoolStats())

		return
	}

	record.TTL, err = dw.getTTL(key)
	if err != nil {

		log.Printf("Error: Get key ttl error, %s\n", err)
		return
	}

	record.DatabaseId = dw.DatabaseId

	dw.writeRecord(record)
	return
}

func (dw *DumpWorker) getSerializeString(key string) (value string, err error) {

	value, err = dw.Client.Dump(key).Result()
	return
}

func (dw *DumpWorker) getTTL(key string) (ttl int64, err error) {

	duration, err := dw.Client.TTL(key).Result()
	ttl = int64(duration.Seconds())
	return
}

func (dw *DumpWorker) writeRecord(record *Record) {

	record.Value = base64.StdEncoding.EncodeToString([]byte(record.Value))
	jsonBytes, err := json.Marshal(record)
	if err != nil {

		log.Printf("Marshal data error , %s\n", err)
		return
	}

	_, err = dw.stream.WriteString(string(jsonBytes) + "\n")
	if err != nil {

		log.Printf("Write file error: %s\n", err)
	}
}

func (dw *DumpWorker) CloseClient() {

	if _, err := dw.Client.Ping().Result(); err != nil {
		return
	}

	dw.Client.Close()
}

func newStream(path string) *os.File {

	fs, err := os.Create(path)
	if err != nil {

		log.Printf("Init file error , %s\n", err)
		return nil
	}

	return fs
}

func Dump(host, password, path string, databaseCount uint64, threadCount int) {

	if databaseCount == 0 {
		databaseCount = getDatabaseCount(host, password)
	}

	stream := newStream(path)
	defer stream.Close()

	var currentDatabase uint64
	for currentDatabase = 0; currentDatabase < databaseCount; currentDatabase++ {

		dumper := &Dumper{
			Client:      createNewClient(host, password, int(currentDatabase), threadCount),
			Host:        host,
			Password:    password,
			DatabaseId:  currentDatabase,
			Stream:      stream,
			ThreadCount: threadCount,
		}
		dumper.Dump()
	}
}

func createNewClient(host, password string, db, poolSize int) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          db,
		PoolSize:    poolSize,
		ReadTimeout: 10,
	})
}

func getDatabaseCount(host, password string) uint64 {
	var databaseCount uint64

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0, // use default DB
	})

	databases, err := client.ConfigGet("databases").Result()
	if err != nil {
		log.Printf("Database config read error, %s\n", err)
		return 0
	}

	if len(databases) == 2 {
		databaseCount, err = strconv.ParseUint(fmt.Sprint(databases[1]), 10, 64)
		if err != nil {

			log.Printf("Read database count error: %s\n", err)
			return 0
		}
	}

	if databaseCount <= 0 {

		log.Printf("Database count read failure\n")
		return 0
	}
	return databaseCount
}
