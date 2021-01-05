package cache

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const maxZipSize = 1024 * 1024 * 10 // 10mb

func copyZipFile(f *zip.File, fs afero.Fs) error {
	if f.CompressedSize64 > maxZipSize {
		return fmt.Errorf("file size %d is larger than maximum %d", f.CompressedSize64, maxZipSize)
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	mf, err := fs.Create(f.Name) // #nosec
	if err != nil {
		return err
	}
	if _, err := io.Copy(mf, io.LimitReader(rc, maxZipSize)); err != nil {
		return nil
	}
	if err := mf.Close(); err != nil {
		return err
	}
	if err := rc.Close(); err != nil {
		return err
	}

	return nil
}

func newFsFromZip(p string) (afero.Fs, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	// Save all contents of zip file to in-memory fs.
	fs := afero.NewMemMapFs()
	for _, f := range reader.File {
		if err := copyZipFile(f, fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func dumpZipFS(fs afero.Fs, path string) error {
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

	return afero.Walk(fs, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := afero.ReadFile(fs, path)
		if err != nil {
			return err
		}

		fw, err := w.Create(path)
		if err != nil {
			return err
		}

		if _, err := fw.Write(data); err != nil {
			return err
		}
		return nil
	})
}
