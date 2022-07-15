package scheduler

import (
	"bytes"
	"fmt"
	"net/url"
)

import (
	"golang.org/x/net/html"
)

/*
get all href in given html node
Params:
	- n: html node
	- refUrl: reference url
*/
func getLinks(n *html.Node, refUrl *url.URL, links []string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				linkUrl, err := refUrl.Parse(a.Val)
				if err == nil {
					links = append(links, linkUrl.String())
				}
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		hl.getLinks(c, refUrl,links)
	}
}
/*
get url links in given html page
Params:
	- data: data for html page
	- urlStr: url string of this html page
Returns:
	- links: parsed links
	- error: any failure
*/
func ParseWebPage(data []byte, urlStr string) ([]string, error) {
	// parse html
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("html.Parse():%s", err.Error())
	}

	// parse url
	refUrl, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, fmt.Errorf("url.ParseRequestURI(%s):%s", urlStr, err.Error())
	}
	
	// get all links
	var links []string
	getLinks(doc, refUrl,links)

	return links, nil
}
