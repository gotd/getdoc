package cache

import (
	"fmt"
)

type ZipReadonly struct {
	files map[string][]byte
}

func NewFromZip(path string) (*ZipReadonly, error) {
	files, err := filesFromZip(path)
	if err != nil {
		return nil, err
	}

	return &ZipReadonly{
		files: files,
	}, nil
}

func (c *ZipReadonly) Get(key string) ([]byte, error) {
	data, found := c.files[key]
	if !found {
		return nil, ErrNotFound
	}

	return data, nil
}

func (c *ZipReadonly) Set(_ string, _ []byte) error {
	return fmt.Errorf("unsupported operation: set")
}
