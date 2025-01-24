package getdoc

import (
	"bytes"
	"context"
	"path"
	"sync"

	"github.com/go-faster/errors"
	"golang.org/x/sync/errgroup"

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

type Extractor struct {
	d Downloader

	layer   int
	workers int
	doc     *Doc
	mux     sync.Mutex
}

func newExtractor(layer, workers int, d Downloader) *Extractor {
	return &Extractor{
		d:       d,
		layer:   layer,
		workers: workers,
	}
}

func (e *Extractor) consume(ctx context.Context, layer int, category, key string) error {
	data, err := e.d.Get(ctx, layer, path.Join(category, key))
	if err != nil {
		return errors.Wrapf(err, "get")
	}

	e.mux.Lock()
	defer e.mux.Unlock()

	reader := bytes.NewReader(data)
	switch category {
	case CategoryConstructor:
		t, err := ParseConstructor(reader)
		if err != nil {
			return errors.Wrap(err, "parse constructor")
		}
		e.doc.Constructors[t.Name] = *t
	case CategoryType:
		t, err := ParseType(reader)
		if err != nil {
			return errors.Wrap(err, "parse type")
		}
		e.doc.Types[t.Name] = *t
	case CategoryMethod:
		t, err := ParseMethod(reader)
		if err != nil {
			return errors.Wrap(err, "parse method")
		}
		e.doc.Methods[t.Name] = *t
		for _, er := range t.Errors {
			e.doc.Errors[er.Type] = er
		}
	}

	return nil
}

func (e *Extractor) Extract(ctx context.Context) (*Doc, error) {
	data, err := e.d.Get(ctx, e.layer, "schema")
	if err != nil {
		return nil, errors.Wrap(err, "get schema")
	}
	index, err := ParseIndex(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "parse index")
	}
	e.doc = &Doc{
		Index:        *index,
		Methods:      map[string]Method{},
		Types:        map[string]Type{},
		Constructors: map[string]Constructor{},
		Errors:       map[string]Error{},
	}

	type Job struct {
		Category, Key string
	}
	jobs := make(chan Job, e.workers)
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < e.workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case job, ok := <-jobs:
					if !ok {
						return nil
					}
					if err := e.consume(ctx, e.layer, job.Category, job.Key); err != nil {
						return errors.Wrap(err, "consume")
					}
				}
			}
		})
	}
	g.Go(func() error {
		defer close(jobs)
		for _, category := range index.Categories {
			for _, v := range category.Values {
				select {
				case jobs <- Job{category.Name, v}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "wait")
	}

	return e.doc, nil
}

// ExtractLayer uses Downloader to extract documentation of specified layer.
func ExtractLayer(ctx context.Context, layer int, d Downloader) (*Doc, error) {
	return newExtractor(layer, 64, d).Extract(ctx)
}
