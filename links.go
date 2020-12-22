package getdoc

import (
	"net/url"
)

func addHost(hrefs []string) (s []string) {
	for _, href := range hrefs {
		u, err := url.Parse(href)
		if err != nil {
			panic(err)
		}

		if u.Host == "" {
			href = "https://core.telegram.org" + href
		}

		s = append(s, href)
	}

	return
}
