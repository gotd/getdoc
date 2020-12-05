package getdoc

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
)

// Constructor represents constructor documentation.
type Constructor struct {
	Name        string
	Description []string
	Fields      map[string]string
}

// ParseConstructor parses html documentation from reader and produces Constructor.
func ParseConstructor(reader io.Reader) (*Constructor, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse document: %w", err)
	}
	return &Constructor{
		Name:        docTitle(doc),
		Description: docDescription(doc),
		Fields:      docParams(doc),
	}, nil
}
