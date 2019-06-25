package acucore

import (
	"golang.org/x/net/html"
)

func checkAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false

}

func checkID(n *html.Node, key string) bool {
	if n.Type == html.ElementNode {
		v, ok := checkAttribute(n, "id")
		if ok && v == key {
			return true
		}
	}
	return false
}

func walk(n *html.Node, key string) *html.Node {
	if checkID(n, key) {
		return n
	}

	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		res := walk(nn, key)
		if res != nil {
			return res
		}
	}
	return nil
}

// GetElementByID query html node by id
func GetElementByID(n *html.Node, key string) *html.Node {
	return walk(n, key)
}
