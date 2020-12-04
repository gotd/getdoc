package getdoc

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIndex(t *testing.T) {
	data, err := ioutil.ReadFile(path.Join("_testdata", "schema.html"))
	if err != nil {
		t.Fatal(err)
	}

	index, err := ParseIndex(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	require.Len(t, index.Categories, 3)
	assert.Len(t, index.Categories[0].Values, 851)
	assert.Len(t, index.Categories[1].Values, 350)
	assert.Len(t, index.Categories[2].Values, 306)
	assert.Equal(t, index.Layer, 121)
}
