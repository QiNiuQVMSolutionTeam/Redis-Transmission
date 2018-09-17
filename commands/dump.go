package commands

import (
	"fmt"
	"github.com/gin-gonic/gin/json"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

type Dumper struct {
	Client *redis.Client
	Path   string
	stream *os.File
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

			d.writeRecord(record)
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	d.Close()
}

func (d *Dumper) Close() {

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

func Dump(host, password, path string) {

	databaseCount := getDatabaseCount(host, password)
	for currentDatabase := 0; uint64(currentDatabase) < databaseCount; currentDatabase++ {

		dumper := &Dumper{
			Client: redis.NewClient(&redis.Options{
				Addr:     host,
				Password: password,        // no password set
				DB:       currentDatabase, // use default DB
			}),
			Path: path,
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
