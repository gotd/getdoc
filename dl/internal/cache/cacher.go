package cache

import "errors"

type Cacher interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

var ErrNotFound = errors.New("not found")
