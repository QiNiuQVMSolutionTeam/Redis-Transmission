package main

import (
	"flag"
	"fmt"
	"github.com/QiNiuQVMSolutionTeam/Redis-Transmission/commands"
	"log"
	"strconv"
)

const ModeDump = "dump"
const ModeRestore = "restore"
const ModeSync = "sync"

func main() {

	var (
		mode                string
		host                string
		password            string
		output              string
		input               string
		databaseCountString string
		sourceHost          string
		destinationHost     string
		sourcePassword      string
		destinationPassword string
	)

	flag.StringVar(&mode, "mode", "", "-mode=[dump|restore]")
	flag.StringVar(&host, "host", "127.0.0.1:6379", "-host=127.0.0.1:6379")
	flag.StringVar(&password, "password", "", "-password=your_password")
	flag.StringVar(&output, "output", "dump.json", "-output=/path/to/file")
	flag.StringVar(&input, "input", "dump.json", "-input=/path/to/file")
	flag.StringVar(&databaseCountString, "database-count", "", "-database-count=16")
	flag.StringVar(&sourceHost, "source", "", "-source=127.0.0.1:6379")
	flag.StringVar(&sourcePassword, "source-password", "", "-source-password=your_password")
	flag.StringVar(&destinationHost, "destination", "", "-destination=127.0.0.1:6378")
	flag.StringVar(&destinationPassword, "destination-password", "", "-destination-password=your_password")

	flag.Parse()

	if mode == ModeDump {

		databaseCount, err := getDatabaseCount(databaseCountString)
		if err != nil {

			log.Printf("Parse database-count err, %s\n", err)
			return
		}

		commands.Dump(host, password, output, databaseCount)

	} else if mode == ModeRestore {

		commands.Restore(host, password, input)

	} else if mode == ModeSync {

		databaseCount, err := getDatabaseCount(databaseCountString)
		if err != nil {

			log.Printf("Parse database-count err, %s\n", err)
			return
		}
		commands.Sync(sourceHost, sourcePassword, destinationHost, destinationPassword, databaseCount)

	} else {

		printHelp()

	}
}

func printHelp() {

	fmt.Print(`
Usage:
	redis-transmission -mode=dump -host=127.0.0.1:6379 [-password=Auth] [-database-count=16] [-output=/path/to/file] [-input=/path/to/file]

	redis-transmission -mode=restore -host=127.0.0.1:6379 [-password=Auth] [-input=/path/to/file]

	redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 [-source-password=Auth] [-destination-password=Auth]

Options:
	-mode=MODE                        Select dump mode, or restore mode. Options: Dump, Restore.
	-host=NODE                        The redis instance (host:port).
	-password=PASSWORD                The redis authorization password, if empty then no use this parameter.
	-input=FILE                       Use for restore data file.
	-output=FILE                      Use for save the dump data file.
	-database-count=COUNT             Specify the redis database count
	-source=NODE                      The source redis instance (host:port).
	-destination=NODE                 The destination redis instance (host:port).
	-source-password=Auth             The source redis authorization password, if empty then no use this parameter.
	-destination-password=Auth        The destination redis authorization password, if empty then no use this parameter.

Examples:
	$ redis-transmission -mode=dump
	$ redis-transmission -mode=dump -host=127.0.0.1:6379
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -output=/tmp/dump.json
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -database-count=16 -output=/tmp/dump.json
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -password=Password -output=/tmp/dump.json
	$ redis-transmission -mode=restore
	$ redis-transmission -mode=restore -host=127.0.0.1:6379
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -input=/tmp/dump.json
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -password=Password -input=/tmp/dump.json
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -password=Password -input=/tmp/dump.json
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -database-count=16
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -source-password=Password -destination-password=Password
`)
}

func getDatabaseCount(databaseCountString string) (databaseCount uint64, err error) {

	if databaseCountString == "" {

		return
	}

	databaseCount, err = strconv.ParseUint(databaseCountString, 10, 64)
	if err != nil {

		return
	}

	return
}
