package href

import (
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// Replace replaces all found HREFs with [index] symbol at end of node text.
// It returns slice of replaced HREFs.
func Replace(selection *goquery.Selection) (hrefs []string) {
	replaceHrefsRecursively(selection, map[int]struct{}{}, &hrefs)
	return hrefs
}

func replaceHrefsRecursively(selection *goquery.Selection, accum map[int]struct{}, hrefs *[]string) {
	if _, processed := accum[selection.Index()]; processed {
		return
	}

	if path, ok := selection.Attr("href"); ok {
		accum[selection.Index()] = struct{}{}

		*hrefs = append(*hrefs, path)

		text, cut := cutRightSpaces(selection.Text())
		text += superscript(len(*hrefs))
		text += cut

		selection.SetText(text)
	}

	selection.Find("*").Each(func(i int, s *goquery.Selection) {
		replaceHrefsRecursively(s, accum, hrefs)
	})
}

func cutRightSpaces(input string) (result, cut string) {
	var (
		r = []rune(input)
		c []rune
	)

	for i := len(r) - 1; i >= 0; i-- {
		if unicode.IsSpace(r[i]) {
			c = append(c, r[i])
			r = r[:i]
		} else {
			break
		}
	}

	return string(r), string(c)
}
