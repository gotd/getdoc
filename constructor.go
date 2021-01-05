package getdoc

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

// Constructor represents constructor documentation.
type Constructor struct {
	Name        string                      `json:"name"`
	Description []string                    `json:"description,omitempty"`
	Links       []string                    `json:"links,omitempty"`
	Fields      map[string]ParamDescription `json:"fields,omitempty"`
}

// ParseConstructor parses html documentation from reader and produces Constructor.
func ParseConstructor(reader io.Reader) (*Constructor, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	desc, links := docDescription(doc)
	return &Constructor{
		Name:        docTitle(doc),
		Description: desc,
		Links:       links,
		Fields:      docParams(doc),
	}, nil
}
