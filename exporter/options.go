package main

import (
	"flag"
	"log"
	"time"
)

type Options struct {
	FetcherEndpoint string
	FetchInterval   time.Duration
	StoreEndpoint   string
	StoreKey        string
	Verbose         bool
}

func NewOptionsFromArgs(_ *log.Logger) *Options {
	fetcherEndpoint := flag.String(
		"fetcher-endpoint",
		"http://localhost:8080",
		"Endpoint used to fetch Jenkins information",
	)

	fetchInterval := flag.Duration(
		"fetch-interval",
		time.Duration(10)*time.Second,
		"Time interval in seconds between fetches",
	)

	storeEndpoint := flag.String(
		"store-endpoint",
		"http://localhost:8889",
		"Endpoint used to store Jenkins information",
	)

	storeKey := flag.String(
		"store-key",
		"key",
		"Key used for storing content",
	)

	verbose := flag.Bool(
		"verbose",
		false,
		"Increase verbosity",
	)

	flag.Parse()

	return &Options{
		FetcherEndpoint: *fetcherEndpoint,
		FetchInterval:   *fetchInterval,
		StoreEndpoint:   *storeEndpoint,
		StoreKey:        *storeKey,
		Verbose:         *verbose,
	}
}
