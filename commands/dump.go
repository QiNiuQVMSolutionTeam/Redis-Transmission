package commands

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

type Dumper struct {
	Client     *redis.Client
	Path       string
	DatabaseId uint64
	stream     *os.File
	Count      uint64
}

func (d *Dumper) Dump() {

	cursor := uint64(0)
	for {
		keys, nextCursor, err := d.scan(cursor)
		if err != nil {

			log.Printf("Error: Scan keys error, %s\n", err)
			break
		}

		for _, key := range keys {

			record := &Record{Key: key}
			record.Value, err = d.getSerializeString(key)
			if err != nil {

				log.Printf("Error: Get key serialize string error, %s\n", err)
				break
			}

			record.TTL, err = d.getTTL(key)
			if err != nil {

				log.Printf("Error: Get key ttl error, %s\n", err)
				break
			}

			record.DatabaseId = d.DatabaseId

			d.writeRecord(record)
			d.Count++

			if d.Count%1000 == 0 {
				d.PrintReport()
			}
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	d.CloseStream()
	d.CloseClient()

	d.PrintReport()
}

func (d *Dumper) CloseStream() {

	if d.stream == nil {
		return
	}

	d.stream.Close()
	d.stream = nil
}

func (d *Dumper) scan(cursor uint64) (keys []string, nextCursor uint64, err error) {

	keys, nextCursor, err = d.Client.Scan(cursor, "", 100).Result()
	return
}

func (d *Dumper) getSerializeString(key string) (value string, err error) {

	value, err = d.Client.Dump(key).Result()
	return
}

func (d *Dumper) getTTL(key string) (ttl int64, err error) {

	duration, err := d.Client.TTL(key).Result()
	ttl = int64(duration.Seconds())
	return
}

func (d *Dumper) writeRecord(record *Record) {

	if !d.initWriter() {

		return
	}

	record.Value = base64.StdEncoding.EncodeToString([]byte(record.Value))
	jsonBytes, err := json.Marshal(record)
	if err != nil {

		log.Printf("Marshal data error , %s\n", err)
		return
	}
	d.stream.Write(jsonBytes)
	d.stream.WriteString("\n")
}

func (d *Dumper) initWriter() bool {

	if d.stream != nil {

		return true
	}

	fs, err := os.Create(d.Path)
	if err != nil {

		log.Printf("Init file error , %s\n", err)
		return false
	}

	d.stream = fs
	return true
}

func (d *Dumper) CloseClient() {

	if _, err := d.Client.Ping().Result(); err != nil {
		return
	}

	d.Client.Close()
}

func (d *Dumper) PrintReport() {

	log.Printf("DB %d dumped %d Record(s).\n", d.DatabaseId, d.Count)
}

func Dump(host, password, path string, databaseCount uint64) {

	if databaseCount == 0 {
		databaseCount = getDatabaseCount(host, password)
	}

	var currentDatabase uint64
	for currentDatabase = 0; currentDatabase < databaseCount; currentDatabase++ {

		dumper := &Dumper{
			Client: redis.NewClient(&redis.Options{
				Addr:     host,
				Password: password,             // no password set
				DB:       int(currentDatabase), // use default DB
			}),
			Path:       path,
			DatabaseId: currentDatabase,
		}
		dumper.Dump()
	}
}

func getDatabaseCount(host, password string) uint64 {
	var databaseCount uint64

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password, // no password set
		DB:       0,        // use default DB
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
