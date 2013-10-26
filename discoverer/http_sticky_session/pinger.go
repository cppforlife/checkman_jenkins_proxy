package http_sticky_session

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Pinger struct {
	endpoint string
	logger   *log.Logger
}

func NewPinger(endpoint string, logger *log.Logger) *Pinger {
	return &Pinger{
		endpoint: endpoint,
		logger:   logger,
	}
}

// Healthy peer will respond with same internal value.
// Returns false when any error occurred.
func (p *Pinger) Ping(peer *Peer) (bool, error) {
	p.logger.Printf("http-sticky-session.pinger.ping endpoint=%s\n", p.endpoint)

	client := p.buildInsecureHttpClient()

	req, err := http.NewRequest("GET", p.endpoint, nil)
	if err != nil {
		p.logger.Printf("http-sticky-session.pinger.ping.build.fail err=%v\n", err)
		return false, err
	}

	req.AddCookie(&http.Cookie{
		Name:     SessionKey,
		Value:    "any-value",
		Path:     "/any-path",
		Domain:   "any-domain",
		Expires:  time.Now().AddDate(0, 0, 1),
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
	})

	req.AddCookie(&http.Cookie{
		Name:     PrivateInstanceKey,
		Value:    peer.LastQueryResult.PrivateInstanceValue,
		Path:     "/any-path",
		Domain:   "any-domain",
		Expires:  time.Now().AddDate(0, 0, 1),
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
	})

	resp, err := client.Do(req)
	if err != nil {
		p.logger.Printf("http-sticky-session.pinger.ping.fail err=%v\n", err)
		return false, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		p.logger.Printf("http-sticky-session.pinger.ping.non-200 err=%v\n", err)
		return false, err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == PrivateInstanceKey {
			return cookie.Value == peer.LastQueryResult.PrivateInstanceValue, nil
		}
	}

	err = errors.New(fmt.Sprintf("%v", resp.Cookies()))
	p.logger.Printf("http-sticky-session.pinger.ping.missing-cookie err=%v\n", err)
	return false, err
}

func (p *Pinger) buildInsecureHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}
