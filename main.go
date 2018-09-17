package main

import (
	"flag"
	"github.com/Luckyboys/RedisDumpRestore/commands"
)

func main() {

	var (
		host     string
		password string
		output   string
	)

	flag.StringVar(&host, "host", "127.0.0.1:6379", "-host=127.0.0.1:6379")
	flag.StringVar(&password, "password", "", "-password=your_password")
	flag.StringVar(&output, "output", "dump.json", "-output=/path/to/file")

	flag.Parse()
	commands.Dump(host, password, output)
}
