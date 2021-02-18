package getdoc

import (
	"encoding/json"
	"errors"

	// For embedding parsed schema.
	_ "embed"
)

// Probably everything below should be code-generated from _schema folder.

// Layers is list of supported layers.
var Layers = []int{
	121,
}

// LayerLatest is id of latest layer.
const LayerLatest = 121

// LayerExists returns true if layer is included in package.
func LayerExists(layer int) bool {
	for _, id := range Layers {
		if id == layer {
			return true
		}
	}
	return false
}

// ErrNotFound means that current package version does not support requested layer.
var ErrNotFound = errors.New("layer not found")

//go:embed _schema/121.json
var schema121 []byte

// Load layer documentation.
func Load(layer int) (*Doc, error) {
	if !LayerExists(layer) {
		return nil, ErrNotFound
	}

	var doc Doc
	if err := json.Unmarshal(schema121, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}
