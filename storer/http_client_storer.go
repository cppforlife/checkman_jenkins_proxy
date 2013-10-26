package storer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpClientStorer struct {
	endpoint string
	logger   *log.Logger
}

func NewHttpClientStorer(endpoint string, logger *log.Logger) *HttpClientStorer {
	return &HttpClientStorer{endpoint, logger}
}

func (hs *HttpClientStorer) Put(key string, content io.Reader) error {
	url := fmt.Sprintf("%s/%s", hs.endpoint, key)
	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, content)
	if err != nil {
		hs.logger.Printf("http-client-storer.put.build.fail err=%v\n", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		hs.logger.Printf("http-client-storer.put.fail err=%v\n", err)
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		hs.logger.Printf("http-client-storer.put.non-200 err=%v\n", err)
		return err
	}

	return nil
}

func (hs *HttpClientStorer) Get(key string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s", hs.endpoint, key)

	resp, err := http.Get(url)
	if err != nil {
		hs.logger.Printf("http-client-storer.get.fail url=%s err=%v\n", url, err)
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		// Body is not closed since return value is ReadCloser
		return resp.Body, nil

	case 404:
		return nil, nil

	default:
		err = errors.New(resp.Status)
		hs.logger.Printf("http-client-storer.get.non-200-or-404 url=%s err=%v\n", url, err)
		return nil, err
	}
}

func (hs *HttpClientStorer) Delete(key string) error {
	url := fmt.Sprintf("%s/%s", hs.endpoint, key)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		hs.logger.Printf("http-client-storer.delete.build.fail err=%v\n", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		hs.logger.Printf("http-client-storer.delete.fail err=%v\n", err)
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		hs.logger.Printf("http-client-storer.delete.non-200 err=%v\n", err)
		return err
	}

	return nil
}
