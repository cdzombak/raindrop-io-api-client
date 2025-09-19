package raindrop

import (
	"golang.org/x/net/html"
	"io"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", false, err
	}

	title, found := traverse(doc)
	return title, found, nil
}
