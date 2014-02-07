package htmlq

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"io"
	"strings"
)

// HtmlQ is a wrapper around a (perhaps empty) set of HTML nodes. It provides
// an interface to both parse and query data.
type HtmlQ struct {
	nodes []*html.Node
}

// ParseString parses the HTML content contained in str. If the markup is
// malformed, it returns an error.
func (hq *HtmlQ) ParseString(str string) error {
	r := strings.NewReader(str)
	return hq.ParseReader(r)
}

// ParseReader consumes all data from the passed reader and parses it. If the 
// reader returns an error, or if the data is malformed, an error is returned.
func (hq *HtmlQ) ParseReader(r io.Reader) error {
	rootNode, er := html.Parse(r)
	if er != nil {
		return er
	}

	hq.nodes = []*html.Node{rootNode}
	return nil
}

// Node returns the underlying *html.Node for the first element in the contained
// nodeset. If this HtmlQ is empty, it returns nil.
func (hq HtmlQ) Node() *html.Node {
	if len(hq.nodes) == 0 {
		return nil
	}

	return hq.nodes[0]
}

// Find returns a new HtmlQ containing only elements in the receiving object
// that match the passed selector.
func (hq HtmlQ) Find(sel string) HtmlQ {
	s, er := cascadia.Compile(sel)
	if er != nil {
		/* XXX: herp */
		panic(er)
	}

	workingSet := []*html.Node{}

	for _, node := range hq.nodes {
		newNodes := s.MatchAll(node)

		for _, newNode := range newNodes {
			alreadyInSet := false

			for _, oldNode := range workingSet {
				if oldNode == newNode {
					alreadyInSet = true
					break
				}
			}

			if !alreadyInSet {
				workingSet = append(workingSet, newNode)
			}
		}
	}

	return HtmlQ{nodes: workingSet}
}

// Val sets the "value" attribute of the first node in the receiving nodeset
// if an argument is given. If no argument is given, it returns the value
// instead.
func (hq HtmlQ) Val(args ...string) string {
	if len(args) == 0 {
		return hq.getVal()
	}

	hq.setVal(args[0])
	return ""
}

func isValueNode(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}

	switch node.Data {
		case "input": fallthrough
		case "textarea": fallthrough
		case "select":
			return true
	}

	return false
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}

func setAttr(node *html.Node, key string, val string) {
	for i, attr := range node.Attr {
		if attr.Key == key {
			node.Attr[i].Val = val
			return
		}
	}

	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: val})
}

func (hq HtmlQ) setVal(val string) {
	for _, node := range hq.nodes {
		if isValueNode(node) {
			setAttr(node, "value", val)
			return
		}
	}
}

func (hq HtmlQ) getVal() string {
	for _, node := range hq.nodes {
		if isValueNode(node) {
			return getAttr(node, "value")
		}
	}

	return ""
}

// Attr is an overloaded function that takes either one or two arguments. In
// its one-argument form, it returns the value of the specified attribute. In
// the two-argument form, it uses the first parameter as the attribute name
// and second as the attribute value. The two-argument form returns nothing.
func (hq HtmlQ) Attr(args ...string) string {
	if len(args) == 0 {
		panic("Attr needs 1 or 2 arguments")
	}

	for _, node := range hq.nodes {
		if node.Type != html.ElementNode {
			continue
		}

		if len(args) == 1 {
			return getAttr(node, args[0])

		} else if len(args) == 2 {
			setAttr(node, args[0], args[1])
			return ""
		}
	}

	return ""
}

// Text returns the concatenated value of all text nodes that descend from
// node N in the nodeset, where N is 0 if no arguments are passed, or the
// value of the first argument passed.
func (hq HtmlQ) Text(args ...int) string {
	idx := 0
	if len(args) > 0 {
		idx = args[0]
	}

	if len(hq.nodes) <= idx {
		return ""
	}

	node := hq.nodes[idx]
	work := stringWriter{}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if er := html.Render(&work, child) ; er != nil {
			/* XXX: derp */
			panic(er)
		}
	}

	return work.String()
}

// Len returns the number of nodes contained within the receiving nodeset.
func (hq HtmlQ) Len() int {
	return len(hq.nodes)
}

// Index returns a new nodeset containing only the node at the specified index.
// If the index is out of bounds, an empty nodeset is returned.
func (hq HtmlQ) Index(idx int) HtmlQ {
	if len(hq.nodes) <= idx {
		return HtmlQ{nodes: nil}
	}

	return HtmlQ{nodes: []*html.Node{hq.nodes[idx]}}
}

// ForEach applies the function f to each node in the receiving nodeset.
func (hq HtmlQ) ForEach(f func(HtmlQ)) HtmlQ {
	for _, node := range hq.nodes {
		f(HtmlQ{nodes:[]*html.Node{node}})
	}

	return hq
}
