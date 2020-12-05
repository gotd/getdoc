package dl

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"path"

	"github.com/cockroachdb/pebble/vfs"
)

func copyZipFile(f *zip.File, fs *vfs.MemFS) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	mf, err := fs.Create(path.Join("zip", f.Name))
	if err != nil {
		return err
	}
	if _, err := io.Copy(mf, io.LimitReader(rc, 1024*1024)); err != nil {
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

func NewZipFS(p string) (vfs.FS, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	// Save all contents of zip file to in-memory fs.
	fs := vfs.NewMem()
	if err := fs.MkdirAll("zip", 0755); err != nil {
		return nil, err
	}
	for _, f := range reader.File {
		if err := copyZipFile(f, fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}
