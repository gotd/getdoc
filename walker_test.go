package getdoc

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/getdoc/dl"
)

type unusableHTTPClient struct{}

func (u unusableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	panic("should not be called")
}

func TestExtract(t *testing.T) {
	fs, err := dl.NewZipFS(filepath.Join("dl", "_testdata", "121.zip"))
	if err != nil {
		t.Fatal(err)
	}

	c, err := dl.NewClient(dl.Options{
		Path:     "zip",
		Client:   unusableHTTPClient{},
		FS:       fs,
		Readonly: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	doc, err := Extract(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 121, doc.Index.Layer)
}
