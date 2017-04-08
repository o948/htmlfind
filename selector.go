package htmlfind

import (
	"regexp"
	"strconv"
	"strings"
)

type selector struct {
	Tag   string
	Attr  string // .class or #id
	Child bool   // >
	Nth   int    // :nth-child() or :nth-last-child()
}

var selectorRx = regexp.MustCompile("^" +
	"(>)?" +
	"([a-z]+)" +
	"([.#][a-zA-Z0-9_-]+)?" +
	"(:(-?[0-9]+))?" +
	"$")

func parseSelectors(s string) (sel []*selector) {
	for _, z := range strings.Fields(s) {
		m := selectorRx.FindStringSubmatch(z)
		if m == nil {
			panic("htmlfind: invalid selector: " + z)
		}
		var r selector
		r.Child = m[1] != ""
		r.Tag = m[2]
		r.Attr = m[3]
		r.Nth, _ = strconv.Atoi(m[5])
		sel = append(sel, &r)
	}
	return
}

func (sel *selector) Matches(node *Node) bool {
	if sel.Tag != node.Tag() {
		return false
	}
	if sel.Attr != "" {
		if sel.Attr[0] == '.' && !hasClass(node.Attr["class"], sel.Attr[1:]) {
			return false
		}
		if sel.Attr[0] == '#' && node.Attr["id"] != sel.Attr[1:] {
			return false
		}
	}
	return true
}

func hasClass(classes, target string) bool {
	for _, class := range strings.Fields(classes) {
		if class == target {
			return true
		}
	}
	return false
}
