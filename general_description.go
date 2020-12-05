package getdoc

import "github.com/PuerkitoBio/goquery"

func generalDescription(doc *goquery.Document) []string {
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
