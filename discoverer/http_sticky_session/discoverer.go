package http_sticky_session

import (
	"fmt"
	"github.com/cppforlife/checkman_jenkins_proxy/discoverer"
	"log"
	"net/http"
	"os"
	"time"
)

type HttpStickySessionDiscoverer struct {
	every        *time.Ticker
	queryReplier *QueryReplier
	querier      *Querier
	pinger       *Pinger
	logger       *log.Logger
}

func NewHttpStickySessionDiscoverer(
	endpoint string,
	every time.Duration,
	logger *log.Logger,
) *HttpStickySessionDiscoverer {

	// Uuid represents a single instance of discoverer
	// which in turn represents a single peer
	uuid := generateUuid()

	return &HttpStickySessionDiscoverer{
		every:        time.NewTicker(every),
		queryReplier: NewQueryReplier(uuid, logger),
		querier:      NewQuerier(endpoint, logger),
		pinger:       NewPinger(endpoint, logger),
		logger:       logger,
	}
}

// Periodically discovers new/existing peers
// and announces them via a callback; blocks
func (hssd *HttpStickySessionDiscoverer) Discover(
	f discoverer.DiscoverCallback) error {

	for {
		select {
		case <-hssd.every.C:
			peer := hssd.queryForPeer()
			if peer != nil {
				f(peer)
			}
		}
	}

	return nil
}

func (hssd *HttpStickySessionDiscoverer) queryForPeer() *Peer {
	queryResult, err := hssd.querier.Query()
	if err != nil {
		hssd.logger.Printf("http-sticky-session-discoverer.discover.query.fail err=%v\n", err)
		return nil
	}

	return NewPeer(queryResult)
}

// net/http.Handler
func (hssd *HttpStickySessionDiscoverer) ServeHTTP(
	respWriter http.ResponseWriter, req *http.Request) {
	hssd.queryReplier.ServeHTTP(respWriter, req)
}

func generateUuid() string {
	file, err := os.Open("/dev/urandom")
	if err != nil {
		panic(fmt.Sprintf("generateUuid: %v", err))
	}

	defer file.Close()

	bytes := make([]byte, 16)
	file.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
