package cache

import (
	"errors"
	"os"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Cacher struct {
	path     string
	files    map[string][]byte
	mux      sync.RWMutex
	readonly bool
}

func NewCacher(path string, readonly bool) (*Cacher, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &Cacher{
				path:     path,
				files:    map[string][]byte{},
				readonly: readonly,
			}, nil
		}

		return nil, err
	}

	files, err := filesFromZip(path)
	if err != nil {
		return nil, err
	}

	return &Cacher{
		path:     path,
		files:    files,
		readonly: readonly,
	}, nil
}

func (c *Cacher) Get(key string) ([]byte, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	data, found := c.files[key]
	if !found {
		return nil, ErrNotFound
	}

	return data, nil
}

func (c *Cacher) Set(key string, value []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.files[key] = value
	return nil
}

func (c *Cacher) Close() error {
	if c.readonly {
		return nil
	}

	return dumpZip(c.files, c.path)
}
