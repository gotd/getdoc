package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type PlainFiles struct {
	path string
}

func NewFromDirectory(path string) PlainFiles {
	return PlainFiles{path}
}

func (p PlainFiles) Get(key string) ([]byte, error) {
	path := filepath.Join(p.path, key)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
	}

	return ioutil.ReadFile(path)
}

func (p PlainFiles) Set(key string, value []byte) error {
	path := filepath.Join(p.path, key)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	return ioutil.WriteFile(path, value, 0600)
}
