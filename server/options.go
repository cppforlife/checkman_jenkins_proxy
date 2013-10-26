package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type Options struct {
	ListenAddr                  string
	HttpDiscovererSelfEndpoint  string
	HttpDiscovererPeersEndpoint string
	HttpDiscovererEvery         time.Duration
	ExpireIn                    time.Duration
	Verbose                     bool
}

func NewOptionsFromArgs(_ *log.Logger) *Options {
	listenAddr := flag.String(
		"listen-address",
		":8889",
		"Address to listen on",
	)

	expireIn := flag.Duration(
		"expire-in",
		time.Duration(30)*time.Second,
		"Time interval in seconds after stored key expires",
	)

	httpDiscovererEndpoint := flag.String(
		"http-discoverer-endpoint",
		"http://checkman_jenkins_proxy.cfapps.io",
		"HTTP discoverer endpoint",
	)

	httpDiscovererEvery := flag.Duration(
		"http-discoverer-every",
		time.Duration(5)*time.Second,
		"Time interval in seconds to discover peers",
	)

	verbose := flag.Bool(
		"verbose",
		false,
		"Increase verbosity",
	)

	flag.Parse()

	return &Options{
		ListenAddr:                  *listenAddr,
		HttpDiscovererSelfEndpoint:  fmt.Sprintf("%s/self", *httpDiscovererEndpoint),
		HttpDiscovererPeersEndpoint: fmt.Sprintf("%s/peers", *httpDiscovererEndpoint),
		HttpDiscovererEvery:         *httpDiscovererEvery,
		ExpireIn:                    *expireIn,
		Verbose:                     *verbose,
	}
}
