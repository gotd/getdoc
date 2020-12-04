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

	var generalDescription []string
	doc.Find("#dev_page_content").Each(func(i int, s *goquery.Selection) {
		s.Children().EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if selection.Is("p") && selection.Text() != "" {
				generalDescription = append(generalDescription, selection.Text())
			}
			return !selection.HasClass("clearfix")
		})
	})

	// 2. Find description of each field.
	fieldDescription := make(map[string]string)
	doc.Find("#dev_page_content > table > tbody > tr").Each(func(i int, row *goquery.Selection) {
		var rowContents []string
		row.Find("td").Each(func(i int, column *goquery.Selection) {
			rowContents = append(rowContents, column.Text())
		})
		if len(rowContents) == 3 {
			fieldDescription[rowContents[0]] = rowContents[2]
		}
	})

	name := doc.Find("#dev_page_title").Text()
	if name == "" {
		return nil, xerrors.New("no title found")
	}

	return &Constructor{
		Name:        name,
		Description: generalDescription,
		Fields:      fieldDescription,
	}, nil
}
