package http_sticky_session

import (
	"encoding/json"
)

type Peer struct {
	LastQueryResult *QueryResult
}

func NewPeer(queryResult *QueryResult) *Peer {
	return &Peer{queryResult}
}

func (p *Peer) Uuid() string {
	return p.LastQueryResult.Uuid
}

// encoding/json JsonMarshaler
func (p *Peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"uuid": p.Uuid(),
	})
}
