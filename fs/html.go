package fs

import (
	"fmt"
	"strings"
)

type node interface {
	Append(n node) node
	AppendText(txt string) node
	AppendTag(txt string) node
	Attr(name string, value string) node
	Render() string
	RenderReadable() string
	renderReadable(indent string) string
}

type htmlNode struct {
	name       string
	txt        string
	children   []node
	attributes []string
}

func createText(txt string) node {
	return &htmlNode{"", txt, []node{}, []string{}}
}

func createNode(name string) node {
	return &htmlNode{name, "", []node{}, []string{}}
}

func (n *htmlNode) setTxt(txt string) node {
	n.txt = txt
	return n
}

func (n *htmlNode) Render() string {
	attrs := n.getAttrs()
	if len(n.txt) > 0 {
		return n.txt
	} else if len(n.children) == 0 && isStandAloneTag(n.name) {
		return fmt.Sprintf("<%s%s>", n.name, attrs)
	}
	txt := n.getInner()
	return fmt.Sprintf("<%s%s>%s</%s>", n.name, attrs, txt, n.name)
}

func (n *htmlNode) RenderReadable() string {
	if len(n.txt) > 0 {
		return n.txt
	}
	txt := n.getInnerReadable("  ")
	attrs := n.getAttrs()
	return fmt.Sprintf("<%s%s>%s\n</%s>",
		n.name, attrs,
		txt,
		n.name)
}

func (n *htmlNode) renderReadable(indent string) string {
	attrs := n.getAttrs()
	if len(n.txt) > 0 {
		return "\n" + indent + n.txt
	} else if len(n.children) == 0 && isStandAloneTag(n.name) {
		return fmt.Sprintf("\n%s<%s%s />",
			indent, n.name, attrs)
	}
	txt := n.getInnerReadable(indent + "  ")
	return fmt.Sprintf("\n%s<%s%s>%s\n%s</%s>",
		indent, n.name, attrs,
		txt,
		indent, n.name)
}

func (n *htmlNode) getInnerReadable(indent string) string {
	txt := ""
	for _, c := range n.children {
		txt += indent + c.renderReadable(indent)
	}
	return txt
}

func (n *htmlNode) getInner() string {
	txt := ""
	for _, c := range n.children {
		txt += c.Render()
	}
	return txt
}

func (n *htmlNode) getAttrs() string {
	if len(n.attributes) > 0 {
		return " " + strings.Join(n.attributes, " ")
	}
	return ""
}

func (n *htmlNode) Append(nd node) node {
	n.children = append(n.children, nd)
	return nd
}

func (n *htmlNode) Attr(name string, value string) node {
	n.attributes = append(
		n.attributes,
		fmt.Sprintf(`%s="%s"`, name, value))
	return n
}

func (n *htmlNode) AppendText(txt string) node {
	n.children = append(n.children, createText(txt))
	return n
}

func (n *htmlNode) AppendTag(nn string) node {
	cn := createNode(nn)
	n.Append(cn)
	return cn
}

type htmlDoc interface {
	Render() string
	AddToHead(n node)
	AddToBody(n node)
}

type html struct {
	doctype string
	head    node
	body    node
}

func newHtmlDoc() htmlDoc {
	return &html{
		"<!doctype html>",
		createNode("head"),
		createNode("body")}
}

func (h *html) Render() string {
	root := createNode("html").Attr("lang", "en")
	root.Append(h.head)
	root.Append(h.body)
	return h.doctype + "\n" + root.Render()
}

func (h *html) AddToHead(n node) {
	h.head.Append(n)
}

func (h *html) AddToBody(n node) {
	h.body.Append(n)
}

func isStandAloneTag(tagname string) bool {
	standalones := []string{"img", "link"}
	for _, t := range standalones {
		if t == tagname {
			return true
		}
	}
	return false
}
