package getdoc

import (
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
)

// Method represents method documentation.
type Method struct {
	Name        string                      `json:"name"`
	Description []string                    `json:"description,omitempty"`
	Links       []string                    `json:"links,omitempty"`
	Parameters  map[string]ParamDescription `json:"parameters,omitempty"`
	Errors      []Error                     `json:"errors,omitempty"`
}

// Error represent possible error documentation.
type Error struct {
	Code        int    `json:"code"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// docErrors extract error code documentation from document.
func docErrors(doc *goquery.Document) []Error {
	var output []Error

	docTableAfter(doc, "#possible-errors").
		Each(func(i int, row *goquery.Selection) {
			var rowContents []string
			row.Find("td").Each(func(i int, column *goquery.Selection) {
				rowContents = append(rowContents, column.Text())
			})
			if len(rowContents) != 3 {
				return
			}
			code, err := strconv.Atoi(rowContents[0])
			if err != nil {
				return
			}
			output = append(output, Error{
				Code:        code,
				Type:        rowContents[1],
				Description: rowContents[2],
			})
		})
	return output
}

// ParseMethod extracts method documentation from reader.
func ParseMethod(reader io.Reader) (*Method, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, errors.Errorf("failed to parse document: %w", err)
	}

	desc, links := docDescription(doc)
	return &Method{
		Name:        docTitle(doc),
		Description: desc,
		Links:       links,
		Parameters:  docParams(doc),
		Errors:      docErrors(doc),
	}, nil
}
