package getdoc

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

// Type represents type (aka class) documentation.
type Type struct {
	Name        string   `json:"name"`
	Description []string `json:"description,omitempty"`
	Links       []string `json:"links,omitempty"`
}

// ParseType parses Type documentation from reader.
func ParseType(reader io.Reader) (*Type, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	desc, links := docDescription(doc)
	return &Type{
		Name:        docTitle(doc),
		Description: desc,
		Links:       links,
	}, nil
}
