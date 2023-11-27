package getdoc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemas(t *testing.T) {
	files, err := data.ReadDir("_schema")
	require.NoError(t, err)
	var latest int
	for _, file := range files {
		layer, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".json"))
		if err != nil {
			continue
		}
		assert.Contains(t, Layers, layer, "layer should be in Layers list")
		if layer > latest {
			latest = layer
		}
	}
	if latest == 0 {
		t.Fatal("no layers found")
	}
	require.Equal(t, latest, LayerLatest, "maximum layer should be latest")
}

func TestLoad(t *testing.T) {
	for _, layer := range append(Layers, LayerLatest) {
		t.Run(fmt.Sprintf("%d", layer), func(t *testing.T) {
			doc, err := Load(layer)
			if err != nil {
				t.Fatal(err)
			}
			if doc.Index.Layer != layer {
				t.Error("layer mismatch")
			}
		})
	}
	t.Run("Latest", func(t *testing.T) {
		doc, err := Load(LayerLatest)
		if err != nil {
			t.Fatal(err)
		}
		if doc.Index.Layer != LayerLatest {
			t.Error("layer mismatch")
		}
	})
	t.Run("NotExist", func(t *testing.T) {
		assert.False(t, LayerExists(-1))
		doc, err := Load(-1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatal("unexpected error")
		}
		require.Nil(t, doc)
	})
}
