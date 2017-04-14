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
	if len(n.txt) > 0 {
		return n.txt
	}
	txt := n.getInner()
	attrs := n.getAttrs()
	return fmt.Sprintf("<%s%s>%s</%s>", n.name, attrs, txt, n.name)
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
	return n
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
