package main

import (
	"chat-app/seeder"
	"flag"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: 'server' or 'seed'")
	flag.Parse()

	if *mode == "seed" {
		seeder.RunSeed()
	} else {
		runServer()
	}
}
