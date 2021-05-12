package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "github.com/KubaiDoLove/gophercises_html_link_parser"
)

type loc struct {
	Value string `xml:"loc"`
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXML := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}
	enc := xml.NewEncoder(os.Stdout)
	fmt.Println(xml.Header)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)
	var q map[string]empty
	nq := map[string]empty{
		urlStr: empty{},
	}

	for {
		q, nq = nq, make(map[string]empty)

		if len(q) == 0 {
			break
		}

		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = empty{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = empty{}
				}
			}

		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, l := range links {
		if keepFn(l) {
			ret = append(ret, l)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
