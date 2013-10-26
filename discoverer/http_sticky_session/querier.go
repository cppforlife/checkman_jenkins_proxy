package http_sticky_session

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SessionKey         = "JSESSIONID"
	PrivateInstanceKey = "__VCAP_ID__"
)

type QueryResult struct {
	Uuid                 string
	PrivateInstanceValue string
}

type Querier struct {
	endpoint string
	logger   *log.Logger
}

func NewQuerier(endpoint string, logger *log.Logger) *Querier {
	return &Querier{
		endpoint: endpoint,
		logger:   logger,
	}
}

// Discovered peer will respond with an internal value.
func (q *Querier) Query() (*QueryResult, error) {
	q.logger.Printf("http-sticky-session.querier.query endpoint=%v\n", q.endpoint)

	client := q.buildInsecureHttpClient()

	req, err := http.NewRequest("GET", q.endpoint, nil)
	if err != nil {
		q.logger.Printf("http-sticky-session.querier.query.build.fail err=%v\n", err)
		return nil, err
	}

	// Adding JSESSIONID triggers CF gorouter to add __VCAP_ID__
	// which forces requests to be forwarded to same instance
	sessionCookie := &http.Cookie{
		Name:     SessionKey,
		Value:    "any-value",
		Path:     "/any-path",
		Domain:   "any-domain",
		Expires:  time.Now().AddDate(0, 0, 1),
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
	}

	req.AddCookie(sessionCookie)

	resp, err := client.Do(req)
	if err != nil {
		q.logger.Printf("http-sticky-session.querier.query.fail err=%v\n", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		q.logger.Printf("http-sticky-session.querier.query.non-200 err=%v\n", err)
		return nil, err
	}

	return q.extractResultFromResponse(resp)
}

func (q *Querier) extractResultFromResponse(resp *http.Response) (*QueryResult, error) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == PrivateInstanceKey {
			uuid, err := q.extractUuidFromResponse(resp)
			if err != nil {
				q.logger.Printf("http-sticky-session.querier.query.extract-uuid err=%v\n", err)
				return nil, err
			}

			return &QueryResult{
				Uuid:                 uuid,
				PrivateInstanceValue: cookie.Value,
			}, nil
		}
	}

	err := errors.New(fmt.Sprintf("%v", resp.Cookies()))
	q.logger.Printf("http-sticky-session.querier.query.missing-cookie err=%v\n", err)
	return nil, err
}

// Returns empty string if error is not nil
func (q *Querier) extractUuidFromResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		q.logger.Printf("http-sticky-session.querier.extract-uuid.read err=%v\n", err)
		return "", err
	}

	var selfJson map[string]string
	err = json.Unmarshal(bytes, &selfJson)
	if err != nil {
		q.logger.Printf("http-sticky-session.querier.extract-uuid.read err=%v\n", err)
		return "", err
	}

	return selfJson["uuid"], nil
}

func (q *Querier) buildInsecureHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}
