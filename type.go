package getdoc

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
)

type Type struct {
	Name        string
	Description []string
}

func ParseType(reader io.Reader) (*Type, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse document: %w", err)
	}

	t := &Type{
		Name:        doc.Find("#dev_page_title").Text(),
		Description: generalDescription(doc),
	}

	return t, nil
}
