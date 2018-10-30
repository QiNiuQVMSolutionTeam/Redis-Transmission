package main

import (
	"flag"
	"fmt"
	"github.com/QiNiuQVMSolutionTeam/Redis-Transmission/commands"
	"log"
	"runtime"
	"strconv"
)

const ModeDump = "dump"
const ModeRestore = "restore"
const ModeSync = "sync"

func main() {

	var (
		mode                          string
		host                          string
		password                      string
		output                        string
		input                         string
		databaseCountString           string
		sourceHost                    string
		destinationHost               string
		sourcePassword                string
		destinationPassword           string
		syncTimesString               string
		threadCountString             string
		isSupportReplaceRestoreString string
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
	flag.StringVar(&syncTimesString, "sync-times", "0", "-sync-times=0")
	flag.StringVar(&threadCountString, "thread-count", strconv.Itoa(runtime.NumCPU()), "-thread-count=4")
	flag.StringVar(&isSupportReplaceRestoreString, "replace-restore", "1", "-replace-restore=1")

	flag.Parse()

	if mode == ModeDump {

		databaseCount, err := getDatabaseCount(databaseCountString)
		if err != nil {

			log.Printf("Parse database-count error, %s\n", err)
			return
		}

		threadCount, err := getThreadCount(threadCountString)
		if err != nil {

			log.Printf("Parse thread-count error, %s\n", err)
			return
		}

		if threadCount <= 0 {

			log.Printf("thread-count parameter error, %s\n", err)
			return
		}

		commands.Dump(host, password, output, databaseCount, threadCount)

	} else if mode == ModeRestore {

		commands.Restore(host, password, input, isSupportReplaceRestoreString != "0")

	} else if mode == ModeSync {

		databaseCount, err := getDatabaseCount(databaseCountString)
		if err != nil {

			log.Printf("Parse database-count err, %s\n", err)
			return
		}
		syncTimes, err := getSyncTimes(syncTimesString)
		if err != nil {

			log.Printf("Parse database-count err, %s\n", err)
			return
		}

		threadCount, err := getThreadCount(threadCountString)
		if err != nil {

			log.Printf("Parse thread-count error, %s\n", err)
			return
		}

		if threadCount <= 0 {

			log.Printf("thread-count parameter error, %s\n", err)
			return
		}

		launcher := &commands.SyncLauncher{}
		launcher.
			SetSourceHost(sourceHost).
			SetSourcePassword(sourcePassword).
			SetDestinationHost(destinationHost).
			SetDestinationPassword(destinationPassword).
			SetDatabaseCount(databaseCount).
			SetSyncTimes(syncTimes).
			SetThreadCount(threadCount).
			SetIsSupportReplaceRestore(isSupportReplaceRestoreString != "0").
			Launch()

	} else {

		printHelp()

	}
}

func printHelp() {

	fmt.Print(`
Usage:
	redis-transmission -mode=dump -host=127.0.0.1:6379 [-password=Auth] [-database-count=16] [-output=/path/to/file] [-input=/path/to/file]

	redis-transmission -mode=restore -host=127.0.0.1:6379 [-password=Auth] [-input=/path/to/file]

	redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 [-source-password=Auth] [-destination-password=Auth] [-database-count=16] [-sync-times=Count]

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
	-sync-times=TIMES                 synchronization times, default loop execution. Do not fill in this parameter if you need to execute it in a loop
	-thread-count=COUNT               Number of concurrent executions, if empty then use cpu cores count.
	-replace-restore=[1|0]            If the destination-side not support restore command use replace option, please use 0 to off this feature, when off this feature, it will remove key before restore command executive, if empty then use replace option.

Examples:
	$ redis-transmission -mode=dump
	$ redis-transmission -mode=dump -host=127.0.0.1:6379
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -output=/tmp/dump.json
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -database-count=16 -output=/tmp/dump.json
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -password=Password -output=/tmp/dump.json
	$ redis-transmission -mode=dump -host=127.0.0.1:6379 -password=Password -output=/tmp/dump.json -thread-count=4
	$ redis-transmission -mode=restore
	$ redis-transmission -mode=restore -host=127.0.0.1:6379
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -input=/tmp/dump.json
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -password=Password -input=/tmp/dump.json
	$ redis-transmission -mode=restore -host=127.0.0.1:6379 -password=Password -input=/tmp/dump.json -replace-restore=0
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -thread-count=16
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -database-count=16
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -source-password=Password -destination-password=Password
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -source-password=Password -destination-password=Password -database-count=1 -sync-times=1
	$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -source-password=Password -destination-password=Password -database-count=1 -replace-restore=0
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

func getThreadCount(threadCountString string) (threadCount int, err error) {

	if threadCountString == "" {

		return
	}

	count, err := strconv.ParseInt(threadCountString, 10, 64)
	if err != nil {

		return
	}

	threadCount = int(count)
	return
}

func getSyncTimes(syncTimesString string) (syncTimes uint64, err error) {

	if syncTimesString == "" {

		return
	}

	syncTimes, err = strconv.ParseUint(syncTimesString, 10, 64)
	if err != nil {

		return
	}

	return
}
