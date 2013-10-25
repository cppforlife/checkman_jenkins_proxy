package main

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
)

type Fetcher interface {
	Fetch() (io.ReadCloser, error)
}

type HttpClientFetcher struct {
	endpoint string
	logger   *log.Logger
}

func NewHttpClientFetcher(endpoint string, logger *log.Logger) *HttpClientFetcher {
	return &HttpClientFetcher{endpoint, logger}
}

func (hcf *HttpClientFetcher) Fetch() (io.ReadCloser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", hcf.endpoint, nil)
	if err != nil {
		hcf.logger.Printf("http-client-fetcher.fetch.build.fail err=%v\n", err)
		return nil, err
	}

	hcf.extractBasicAuthFromURL(req)

	resp, err := client.Do(req)
	if err != nil {
		hcf.logger.Printf("http-client-fetcher.fetch.fail err=%v\n", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		hcf.logger.Printf("http-client-fetcher.fetch.fail.non-200 err=%v\n", err)
		return nil, err
	}

	// Body is not closed since return value is ReadCloser
	return resp.Body, nil
}

func (hcf *HttpClientFetcher) extractBasicAuthFromURL(req *http.Request) {
	if req.URL.User != nil {
		username := req.URL.User.Username()
		password, _ := req.URL.User.Password()
		req.SetBasicAuth(username, password)
	}

	req.URL.User = nil
}
