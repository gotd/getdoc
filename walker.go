package getdoc

import (
	"bytes"
	"context"
	"path"

	"golang.org/x/xerrors"
)

type Downloader interface {
	Get(ctx context.Context, layer int, key string) ([]byte, error)
}

type walker struct {
	client Downloader
}

func (w *walker) Index(ctx context.Context) (*Index, error) {
	data, err := w.client.Get(ctx, 0, "schema")
	if err != nil {
		return nil, err
	}
	return ParseIndex(bytes.NewReader(data))
}

// Extracts uses Downloader to extract documentation.
func Extract(ctx context.Context, d Downloader) (*Doc, error) {
	w := walker{
		client: d,
	}

	index, err := w.Index(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to extract index: %w", err)
	}
	doc := &Doc{
		Index:        *index,
		Methods:      map[string]Method{},
		Types:        map[string]Type{},
		Constructors: map[string]Constructor{},
	}
	for _, category := range index.Categories {
		for _, v := range category.Values {
			data, err := w.client.Get(ctx, index.Layer, path.Join(category.Name, v))
			if err != nil {
				return nil, xerrors.Errorf("fetch(%s/%s) failed: %w", category.Name, v, err)
			}
			reader := bytes.NewReader(data)
			switch category.Name {
			case CategoryConstructor:
				t, err := ParseConstructor(reader)
				if err != nil {
					return nil, xerrors.Errorf("parse(%s/%) failed: %w", category.Name, v, err)
				}
				doc.Constructors[t.Name] = *t
			case CategoryType:
				t, err := ParseType(reader)
				if err != nil {
					return nil, xerrors.Errorf("parse(%s/%) failed: %w", category.Name, v, err)
				}
				doc.Types[t.Name] = *t
			case CategoryMethod:
				t, err := ParseMethod(reader)
				if err != nil {
					return nil, xerrors.Errorf("parse(%s/%) failed: %w", category.Name, v, err)
				}
				doc.Methods[t.Name] = *t
			}
		}
	}

	return doc, nil
}
