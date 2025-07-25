package getdoc

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"

	"github.com/go-faster/errors"
)

// Probably everything below should be code-generated from _schema folder.

// Layers is list of supported layers.
var Layers = []int{
	121,
	133,
	138,
	139,
	144,
	145,
	155,
	158,
	164,
	167,
	170,
	181,
	185,
	192,
	195,
	208,
}

// LayerLatest is id of the latest layer.
const LayerLatest = 208

// LayerExists returns true if layer is included in package.
func LayerExists(layer int) bool {
	for _, v := range Layers {
		if v == layer {
			return true
		}
	}
	return false
}

// ErrNotFound means that current package version does not support requested layer.
var ErrNotFound = errors.New("layer not found")

//go:embed _schema
var embedData embed.FS

// Load layer documentation.
func Load(layer int) (*Doc, error) {
	if !LayerExists(layer) {
		return nil, ErrNotFound
	}

	b, err := embedData.ReadFile(path.Join("_schema", fmt.Sprintf("%d.json", layer)))
	if err != nil {
		return nil, errors.Wrap(err, "read")
	}

	var doc Doc
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, errors.Wrap(err, "parse")
	}

	return &doc, nil
}
