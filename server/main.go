package main

import (
	discoverer_ "github.com/cppforlife/checkman_jenkins_proxy/discoverer"
	"github.com/cppforlife/checkman_jenkins_proxy/discoverer/http_sticky_session"
	"github.com/cppforlife/checkman_jenkins_proxy/storer"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "server: ", log.Ltime)

	options := NewOptionsFromArgs(logger)
	logger.Printf("options=%v\n", options)

	memoryStorer := storer.NewMemoryStorer(logger)

	expiringStorer := storer.NewExpiringStorer(
		memoryStorer, options.ExpireIn, logger)

	discoverer := http_sticky_session.NewHttpStickySessionDiscoverer(
		options.HttpDiscovererSelfEndpoint,
		options.HttpDiscovererEvery,
		logger,
	)

	discovererSelfUrl, _ := url.Parse(options.HttpDiscovererSelfEndpoint)
	discovererPeersUrl, _ := url.Parse(options.HttpDiscovererPeersEndpoint)

	peerCollection := discoverer_.NewPeerCollection(logger)

	server := NewServer(
		options.ListenAddr,
		expiringStorer,
		map[string]http.Handler{
			discovererSelfUrl.Path:  discoverer,
			discovererPeersUrl.Path: peerCollection,
		},
		logger,
	)

	go discoverer.Discover(func(peer discoverer_.Peer) {
		peerCollection.AddOrUpdateByUuid(peer)
	})

	server.ListenAndServe()
}
