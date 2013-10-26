package discoverer

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type PeerCollection struct {
	sync.RWMutex
	peers  *list.List
	logger *log.Logger
}

func NewPeerCollection(logger *log.Logger) *PeerCollection {
	return &PeerCollection{
		peers:  list.New(),
		logger: logger,
	}
}

func (pc *PeerCollection) AddOrUpdateByUuid(peer Peer) {
	pc.Lock()
	defer pc.Unlock()

	pc.logger.Printf("peer-collection.add-or-update-by-unique-id uuid=%s\n", peer.Uuid())
	for el := pc.peers.Front(); el != nil; el = el.Next() {
		if el.Value.(Peer).Uuid() == peer.Uuid() {
			pc.peers.Remove(el)
			break
		}
	}

	pc.peers.PushBack(peer)
}

// net/http.Handler
func (pc *PeerCollection) ServeHTTP(
	respWriter http.ResponseWriter, req *http.Request) {

	bytes, err := json.Marshal(pc)
	if err != nil {
		pc.logger.Printf("discoverer.peer-collection.serve-http.fail err=%v\n", err)
		pc.respondWithError(500, respWriter)
	}

	respWriter.Write(bytes)
}

func (pc *PeerCollection) respondWithError(code int, respWriter http.ResponseWriter) {
	pc.logger.Printf("discoverer.peer-collection.respond-with-error code=%d\n", code)
	body := fmt.Sprintf("%d %s", code, http.StatusText(code))
	http.Error(respWriter, body, code)
}

// encoding/json.Marshaler
func (pc *PeerCollection) MarshalJSON() ([]byte, error) {
	pc.Lock()
	defer pc.Unlock()

	var i int = 0
	var peers []Peer = make([]Peer, pc.peers.Len())

	for el := pc.peers.Front(); el != nil; el = el.Next() {
		peers[i] = el.Value.(Peer)
		i++
	}
	return json.Marshal(peers)
}
