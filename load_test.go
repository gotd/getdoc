package getdoc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	for _, layer := range append(Layers) {
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
