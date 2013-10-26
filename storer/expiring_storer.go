package storer

import (
	"io"
	"log"
	"sync"
	"time"
)

type ExpiringStorer struct {
	sync.RWMutex
	storer   Storer
	expireIn time.Duration
	times    map[string]time.Time
	logger   *log.Logger
}

func NewExpiringStorer(
	storer Storer,
	expireIn time.Duration,
	logger *log.Logger,
) *ExpiringStorer {

	return &ExpiringStorer{
		storer:   storer,
		expireIn: expireIn,
		times:    make(map[string]time.Time),
		logger:   logger,
	}
}

func (es *ExpiringStorer) Put(key string, content io.Reader) error {
	es.Lock()
	es.times[key] = time.Now()
	es.Unlock()

	return es.storer.Put(key, content)
}

func (es *ExpiringStorer) Get(key string) (io.ReadCloser, error) {
	es.Lock()
	putAt, hasKey := es.times[key]
	es.Unlock()

	content, err := es.storer.Get(key)
	if err == nil {
		if content != nil && !hasKey {
			es.logger.Panicf("expiring-storer.get.missing-put-at key=%s\n", key)
		}

		if time.Since(putAt) > es.expireIn {
			es.logger.Printf("expiring-storer.get.expire key=%s\n", key)

			err := es.storer.Delete(key)
			if err != nil {
				es.logger.Printf("expiring-storer.get.delete.fail key=%s err=%v\n", key, err)
				// ignore since deletion is a hidden operation in this case
			}

			return nil, nil
		}
	}

	return content, err
}

func (es *ExpiringStorer) Delete(key string) error {
	es.Lock()
	delete(es.times, key)
	es.Unlock()

	return es.storer.Delete(key)
}
