package main

import (
	"github.com/cppforlife/checkman_jenkins_proxy/storer"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "exporter: ", log.Ltime)

	options := NewOptionsFromArgs(logger)
	logger.Printf("options=%v\n", options)

	exporter := NewPeriodicExporter(
		NewPeriodicFetchScheduler(options.FetchInterval, logger),
		NewHttpClientFetcher(options.FetcherEndpoint, logger),
		storer.NewHttpClientStorer(options.StoreEndpoint, logger),
		options.StoreKey,
		logger,
	)

	exporter.Run()
}
