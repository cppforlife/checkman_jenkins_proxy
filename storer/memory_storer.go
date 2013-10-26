package storer

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"sync"
)

type MemoryStorer struct {
	sync.RWMutex
	memory map[string][]byte
	logger *log.Logger
}

func NewMemoryStorer(logger *log.Logger) *MemoryStorer {
	return &MemoryStorer{
		memory: make(map[string][]byte),
		logger: logger,
	}
}

func (ms *MemoryStorer) Put(key string, content io.Reader) error {
	ms.Lock()
	defer ms.Unlock()

	bytes, err := ioutil.ReadAll(content)
	if err != nil {
		ms.logger.Printf("memory-storer.put.fail err=%v\n", err)
		return err
	}

	ms.memory[key] = bytes
	ms.logger.Printf("memory-storer.put.success key=%s len=%d\n", key, len(bytes))

	return nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func (ms *MemoryStorer) Get(key string) (io.ReadCloser, error) {
	ms.Lock()
	defer ms.Unlock()

	content := ms.memory[key]
	if content == nil {
		ms.logger.Printf("memory-storer.get.nil key=%s\n", key)
		return nil, nil
	}

	ms.logger.Printf("memory-storer.get.success key=%s\n", key)
	return nopCloser{bytes.NewBuffer(content)}, nil
}

func (ms *MemoryStorer) Delete(key string) error {
	ms.Lock()
	defer ms.Unlock()

	delete(ms.memory, key)

	ms.logger.Printf("memory-storer.delete.success key=%s\n", key)
	return nil
}
