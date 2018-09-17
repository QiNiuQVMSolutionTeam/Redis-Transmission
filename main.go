package main

import (
	"flag"
	"fmt"
	"github.com/Luckyboys/RedisDumpRestore/commands"
)

const ModeDump = "dump"
const ModeRestore = "restore"

func main() {

	var (
		mode     string
		host     string
		password string
		output   string
		input    string
	)

	flag.StringVar(&mode, "mode", "", "-mode=[dump|restore]")
	flag.StringVar(&host, "host", "127.0.0.1:6379", "-host=127.0.0.1:6379")
	flag.StringVar(&password, "password", "", "-password=your_password")
	flag.StringVar(&output, "output", "dump.json", "-output=/path/to/file")
	flag.StringVar(&input, "input", "dump.json", "-input=/path/to/file")

	flag.Parse()

	if mode == ModeDump {

		commands.Dump(host, password, output)

	} else if mode == ModeRestore {

		commands.Restore(host, password, input)

	} else {

		printHelp()

	}
}

func printHelp() {

	fmt.Print(`
Usage:
	redis-dump-restore -mode=[dump|restore] -host=127.0.0.1:6379 [-password=Auth] [-output=/path/to/file] [-input=/path/to/file]

Options:
	-mode=MODE                        Select dump mode, or restore mode. Options: Dump, Restore.
	-host=NODE                        The redis instance (host:port).
	-password=PASSWORD                The redis authorization password, if empty then no use this parameter.
	-input=FILE                       Use for restore data file.
	-output=FILE                      Use for save the dump data file.

Examples:
	$ redis-dump-restore -mode=dump
	$ redis-dump-restore -mode=dump -host=127.0.0.1:6379
	$ redis-dump-restore -mode=dump -host=127.0.0.1:6379 -output=/tmp/dump.json
	$ redis-dump-restore -mode=dump -host=127.0.0.1:6379 -password=Password -output=/tmp/dump.json
	$ redis-dump-restore -mode=restore
	$ redis-dump-restore -mode=restore -host=127.0.0.1:6379
	$ redis-dump-restore -mode=restore -host=127.0.0.1:6379 -input=/tmp/dump.json
	$ redis-dump-restore -mode=restore -host=127.0.0.1:6379 -password=Password -input=/tmp/dump.json
`)
}
