package cache

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
)

func filesFromZip(p string) (map[string][]byte, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, f := range reader.File {
		if err := func() error {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer func() { _ = rc.Close() }()

			buf, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			files[f.Name] = buf
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	return files, nil
}
