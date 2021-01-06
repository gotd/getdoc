package getdoc

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Category index.
type Category struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// Index of schema.
type Index struct {
	Layer      int        `json:"layer"`
	Categories []Category `json:"categories"`
}

type indexLink struct {
	Type  string
	Value string
}

// ParseIndex parses schema index documentation from reader.
func ParseIndex(reader io.Reader) (*Index, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Searching for current layer.
	var layer int
	doc.Find("a.dropdown-toggle").Each(func(i int, selection *goquery.Selection) {
		matches := regexp.MustCompile(`Layer (\d+)`).FindStringSubmatch(selection.Text())
		id, err := strconv.Atoi(matches[1])
		if err == nil {
			layer = id
		}
	})
	if layer == 0 {
		return nil, errors.New("unable to find layer id")
	}

	// Searching for all unique references.
	links := map[indexLink]struct{}{}
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

	// Aggregating references per category.
	categories := map[string][]string{}
	for l := range links {
		categories[l.Type] = append(categories[l.Type], l.Value)
	}

	// Sorting categories.
	var categoryNames []string
	for name := range categories {
		categoryNames = append(categoryNames, name)
	}
	sort.Strings(categoryNames)

	// Aggregating back.
	index := &Index{
		Layer: layer,
	}
	for _, name := range categoryNames {
		values := categories[name]
		sort.Strings(values)
		index.Categories = append(index.Categories, Category{
			Name:   name,
			Values: values,
		})
	}

	return index, nil
}

const (
	CategoryType        = "type"
	CategoryConstructor = "constructor"
	CategoryMethod      = "method"
)

func extractLink(v string) indexLink {
	for _, s := range []string{
		CategoryType,
		CategoryConstructor,
		CategoryMethod,
	} {
		pref := "/" + s + "/"
		if strings.HasPrefix(v, pref) {
			val := strings.TrimPrefix(v, pref)
			val = strings.TrimSpace(val)
			unescaped, err := url.PathUnescape(val)
			if err == nil {
				val = unescaped
			}
			return indexLink{
				Type:  s,
				Value: val,
			}
		}
	}
	return indexLink{}
}
