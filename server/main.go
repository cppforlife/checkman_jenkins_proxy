package main

import (
	"github.com/cppforlife/checkman_jenkins_proxy/server/storer"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "server: ", log.Ltime)

	options := NewOptionsFromArgs(logger)
	logger.Printf("options=%v\n", options)

	memoryStorer := storer.NewMemoryStorer(logger)

	expiringStorer := storer.NewExpiringStorer(
		memoryStorer, options.ExpireIn, logger)

	server := NewServer(options.ListenAddr, expiringStorer, logger)

	server.ListenAndServe()
}
