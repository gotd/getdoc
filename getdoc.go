// Package getdoc provides a way to transform Telegram TL documentation into
// machine-readable format.
package getdoc

import (
	"github.com/PuerkitoBio/goquery"
)

// Doc represents full documentation description.
type Doc struct {
	Index Index `json:"index"`

	Constructors map[string]Constructor `json:"constructors"`
	Types        map[string]Type        `json:"types"`
	Methods      map[string]Method      `json:"methods"`
}

// docTitle extracts title from document.
func docTitle(doc *goquery.Document) string {
	return doc.Find("#dev_page_title").Text()
}

// docDescription extracts description lines from document.
func docDescription(doc *goquery.Document) []string {
	var description []string
	doc.Find("#dev_page_content").Each(func(i int, s *goquery.Selection) {
		s.Children().EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if selection.Is("p") && selection.Text() != "" {
				description = append(description, selection.Text())
			}
			return !selection.HasClass("clearfix")
		})
	})
	return description
}

// docTableAfter extracts table after selector "after".
func docTableAfter(doc *goquery.Document, after string) *goquery.Selection {
	var (
		meetAfter bool
		table     *goquery.Selection
	)
	doc.Find("#dev_page_content").Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find(after).Length() > 0 {
			// Found title of table. Next <table> element will be requested table.
			meetAfter = true
			return true
		}
		if meetAfter && s.Is("table") {
			// Found requested table, stopping iteration.
			table = s
			return false
		}
		return true
	})
	if table == nil {
		return &goquery.Selection{}
	}
	return table.First().Find("tbody > tr")
}

// docParams extract parameters documentation from document.
//
// Key is parameter name, value is documentation string.
func docParams(doc *goquery.Document) map[string]string {
	fields := make(map[string]string)
	docTableAfter(doc, "#parameters").
		Each(func(i int, row *goquery.Selection) {
			var rowContents []string
			row.Find("td").Each(func(i int, column *goquery.Selection) {
				rowContents = append(rowContents, column.Text())
			})
			if len(rowContents) == 3 {
				fields[rowContents[0]] = rowContents[2]
			}
		})
	return fields
}
