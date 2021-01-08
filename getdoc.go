// Package getdoc provides a way to transform Telegram TL documentation into
// machine-readable format.
package getdoc

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/gotd/getdoc/href"
)

// Doc represents full documentation description.
type Doc struct {
	Index Index `json:"index"`

	Constructors map[string]Constructor `json:"constructors"`
	Types        map[string]Type        `json:"types"`
	Methods      map[string]Method      `json:"methods"`
	Errors       map[string]Error       `json:"errors"`
}

// docTitle extracts title from document.
func docTitle(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("#dev_page_title").Text())
}

// docDescription extracts description lines from document.
func docDescription(doc *goquery.Document) (desc, links []string) {
	doc.Find("#dev_page_content").Each(func(i int, s *goquery.Selection) {
		s.Children().EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if selection.Is("p") && selection.Text() != "" {
				hrefs := href.Replace(selection)

				text := strings.TrimSpace(selection.Text())
				for _, part := range strings.Split(text, "\n") {
					part = strings.TrimSpace(part)
					if part == "" {
						continue
					}
					desc = append(desc, part)
				}

				links = append(links, addHost(hrefs)...)
			}
			return !selection.HasClass("clearfix")
		})
	})
	return
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

type ParamDescription struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Links       []string `json:"links"`
}

// docParams extract parameters documentation from document.
//
// Key is parameter name, value is documentation struct.
func docParams(doc *goquery.Document) map[string]ParamDescription {
	fields := make(map[string]ParamDescription)

	docTableAfter(doc, "#parameters").
		Each(func(i int, row *goquery.Selection) {
			var rowContents []string
			var links []string
			row.Find("td").Each(func(i int, column *goquery.Selection) {
				links = addHost(href.Replace(column))
				rowContents = append(rowContents, column.Text())
			})
			if len(rowContents) == 3 {
				fields[rowContents[0]] = ParamDescription{
					Name:        rowContents[0],
					Description: rowContents[2],
					Links:       links,
				}
			}
		})
	return fields
}
