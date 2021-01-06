package getdoc

import (
	"bytes"
	"context"
	"fmt"
	"path"

	"github.com/gotd/getdoc/dl"
)

// Downloader abstracts documentation fetching.
type Downloader interface {
	Get(ctx context.Context, layer int, key string) ([]byte, error)
}

// Extract uses Downloader to extract documentation.
func Extract(ctx context.Context, d Downloader) (*Doc, error) {
	return ExtractLayer(ctx, dl.NoLayer, d)
}

// ExtractLayer uses Downloader to extract documentation of specified layer.
func ExtractLayer(ctx context.Context, layer int, d Downloader) (*Doc, error) {
	data, err := d.Get(ctx, layer, "schema")
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}
	index, err := ParseIndex(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to extract index: %w", err)
	}
	doc := &Doc{
		Index:        *index,
		Methods:      map[string]Method{},
		Types:        map[string]Type{},
		Constructors: map[string]Constructor{},
	}
	for _, category := range index.Categories {
		for _, v := range category.Values {
			data, err := d.Get(ctx, index.Layer, path.Join(category.Name, v))
			if err != nil {
				return nil, fmt.Errorf("fetch(%s/%s) failed: %w", category.Name, v, err)
			}
			reader := bytes.NewReader(data)
			switch category.Name {
			case CategoryConstructor:
				t, err := ParseConstructor(reader)
				if err != nil {
					return nil, fmt.Errorf("parse(%s/%s) failed: %w", category.Name, v, err)
				}
				doc.Constructors[t.Name] = *t
			case CategoryType:
				t, err := ParseType(reader)
				if err != nil {
					return nil, fmt.Errorf("parse(%s/%s) failed: %w", category.Name, v, err)
				}
				doc.Types[t.Name] = *t
			case CategoryMethod:
				t, err := ParseMethod(reader)
				if err != nil {
					return nil, fmt.Errorf("parse(%s/%s) failed: %w", category.Name, v, err)
				}
				doc.Methods[t.Name] = *t
			}
		}
	}

	return doc, nil
}
