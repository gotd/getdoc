package getdoc

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
)

// Constructor represents constructor documentation.
type Constructor struct {
	Name        string            `json:"name"`
	Description []string          `json:"description,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
}

// ParseConstructor parses html documentation from reader and produces Constructor.
func ParseConstructor(reader io.Reader) (*Constructor, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, errors.Errorf("failed to parse document: %w", err)
	}
	return &Constructor{
		Name:        docTitle(doc),
		Description: docDescription(doc),
		Fields:      docParams(doc),
	}, nil
}
