package cache

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func filesFromZip(p string) (map[string][]byte, error) {
	data, err := ioutil.ReadFile(p)
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

			buf, err := ioutil.ReadAll(rc)
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

func dumpZip(files map[string][]byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for path, data := range files {
		path = filepath.ToSlash(path)
		fw, err := w.Create(path)
		if err != nil {
			return err
		}

		if _, err := fw.Write(data); err != nil {
			return err
		}
	}

	return nil
}
