package htmlfind

import (
	"io"

	"golang.org/x/net/html"
)

type Node struct {
	node     *html.Node
	Attr     map[string]string
	Children []*Node
}

func newNode(node *html.Node) *Node {
	attr := make(map[string]string)
	for _, a := range node.Attr {
		attr[a.Key] = a.Val
	}
	var children []*Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			children = append(children, newNode(child))
		}
	}
	return &Node{node, attr, children}
}

func Parse(r io.Reader) (*Node, error) {
	if node, err := html.Parse(r); err != nil {
		return nil, err
	} else {
		return newNode(node), nil
	}
}

func (root *Node) Tag() string {
	return root.node.Data
}

func (root *Node) Text() string {
	return text(root.node)
}

func (root *Node) FindN(query string, n int) []*Node {
	return root.find(parseSelectors(query), nil, n)
}

func (root *Node) FindFirst(query string) *Node {
	if found := root.find(parseSelectors(query), nil, 1); found != nil {
		return found[0]
	}
	return nil
}

func (root *Node) find(sel []*selector, found []*Node, n int) []*Node {
	if len(found) == n {
		return found
	}
	if len(sel) == 0 {
		return append(found, root)
	}

	if sel[0].Nth == 0 {
		for _, child := range root.Children {
			if sel[0].Matches(child) {
				found = child.find(sel[1:], found, n)
			}
		}
	} else {
		var i int
		if k := sel[0].Nth; k > 0 {
			i = k - 1
		} else {
			k = -k
			i = len(root.Children) - k
		}
		if 0 <= i && i < len(root.Children) {
			child := root.Children[i]
			if sel[0].Matches(child) {
				found = child.find(sel[1:], found, n)
			}
		}
	}

	if !sel[0].Child {
		for _, child := range root.Children {
			found = child.find(sel, found, n)
		}
	}

	return found
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var buf string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		buf += text(child)
	}
	return buf
}
