package storer

import (
	"io"
)

type Storer interface {
	Put(key string, content io.Reader) error
	Get(key string) (io.ReadCloser, error)
	Delete(key string) error
}
