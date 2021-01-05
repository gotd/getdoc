package cache

import (
	"errors"
	"os"

	"github.com/spf13/afero"
)

var ErrNotFound = errors.New("not found")

type Cacher struct {
	path     string
	fs       afero.Fs
	readonly bool
}

func NewCacher(path string, readonly bool) (*Cacher, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &Cacher{
				path:     path,
				fs:       afero.NewMemMapFs(),
				readonly: readonly,
			}, nil
		}

		return nil, err
	}

	fs, err := newFsFromZip(path)
	if err != nil {
		return nil, err
	}

	return &Cacher{
		path:     path,
		fs:       fs,
		readonly: readonly,
	}, nil
}

func (c *Cacher) Get(key string) ([]byte, error) {
	if _, err := c.fs.Stat(key); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return afero.ReadFile(c.fs, key)
}

func (c *Cacher) Set(key string, value []byte) error {
	return afero.WriteFile(c.fs, key, value, 0600)
}

func (c *Cacher) Close() error {
	if c.readonly {
		return nil
	}

	return dumpZipFS(c.fs, c.path)
}
