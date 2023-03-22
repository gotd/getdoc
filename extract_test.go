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

func (u unusableHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	panic("should not be called")
}

func TestExtract(t *testing.T) {
	c, err := dl.NewClient(dl.Options{
		Client:   unusableHTTPClient{},
		Path:     filepath.Join("dl", "_testdata", "121.zip"),
		Readonly: true,
		FromZip:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	doc, err := ExtractLayer(context.Background(), 121, c)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, 121, doc.Index.Layer)
}
