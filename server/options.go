package main

import (
	"flag"
	"log"
	"time"
)

type Options struct {
	ListenAddr string
	ExpireIn   time.Duration
	Verbose    bool
}

func NewOptionsFromArgs(_ *log.Logger) *Options {
	listenAddr := flag.String(
		"listen-address", ":8889", "Address to listen on")

	expireIn := flag.Duration(
		"expire-in",
		time.Duration(30)*time.Second,
		"Time interval in seconds after stored key expires",
	)

	verbose := flag.Bool(
		"verbose", false, "Increase verbosity")

	flag.Parse()

	return &Options{
		ListenAddr: *listenAddr,
		ExpireIn:   *expireIn,
		Verbose:    *verbose,
	}
}
