package href

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestHref(t *testing.T) {
	tests := []struct {
		Input string
		Text  string
		HREFs []string
	}{
		{
			Input: `<p>Hello! <a href="https://foo.com/bar">Click me!</a></p>`,
			Text:  `Hello! Click me!¹`,
			HREFs: []string{"https://foo.com/bar"},
		},
		{
			Input: `<p>Hello! <a href="/foo">Click me</a>   again!</p>`,
			Text:  `Hello! Click me¹   again!`,
			HREFs: []string{"/foo"},
		},
	}

	for _, test := range tests {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(test.Input))
		assert.NoError(t, err)
		hrefs := Replace(doc.Selection)
		assert.Equal(t, test.Text, doc.Text())
		assert.Equal(t, hrefs, test.HREFs)
	}
}
