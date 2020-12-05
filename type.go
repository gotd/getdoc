package getdoc

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
)

// Type represents type (aka class) documentation.
type Type struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
}

// ParseType parses Type documentation from reader.
func ParseType(reader io.Reader) (*Type, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, errors.Errorf("failed to parse document: %w", err)
	}
	return &Type{
		Name:        docTitle(doc),
		Description: docDescription(doc),
	}, nil
}
