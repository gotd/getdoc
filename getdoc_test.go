package getdoc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestIndex(t *testing.T) {
	data, err := ioutil.ReadFile(path.Join("_testdata", "schema.html"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	links := map[Link]struct{}{}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		v, ok := selection.Attr("href")
		if !ok {
			return
		}
		l := extractLink(v)
		if l.Type == "" {
			return
		}
		links[l] = struct{}{}
	})

	fmt.Println("links", len(links))
}

type Category struct {
	Name   string
	Values []string
}

type Link struct {
	Type  string
	Value string
}

func extractLink(v string) Link {
	for _, s := range []string{
		"type",
		"constructor",
		"method",
	} {
		pref := "/" + s + "/"
		if strings.HasPrefix(v, "/"+s+"/") {
			val := strings.TrimPrefix(v, pref)
			val = strings.TrimSpace(val)
			unescaped, err := url.PathUnescape(val)
			if err == nil {
				val = unescaped
			}
			return Link{
				Type:  s,
				Value: val,
			}
		}
	}
	return Link{}
}
