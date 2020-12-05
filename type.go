package getdoc

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
)

// Type represents type (aka class) documentation.
type Type struct {
	Name        string
	Description []string
}

// ParseType parses Type documentation from reader.
func ParseType(reader io.Reader) (*Type, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse document: %w", err)
	}
	return &Type{
		Name:        docTitle(doc),
		Description: docDescription(doc),
	}, nil
}
