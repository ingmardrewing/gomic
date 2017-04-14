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
	standalones := []string{"img", "link", "meta"}
	for _, t := range standalones {
		if t == tagname {
			return true
		}
	}
	return false
}

type htmlDocWrapperI interface {
	Render() string
	Init()
}

type htmlDocWrapper struct {
	htmlDoc htmlDoc
}

func newHtmlDocWrapper() htmlDocWrapperI {
	return &htmlDocWrapper{newHtmlDoc()}
}

func (hdw *htmlDocWrapper) Init() {
	hdw.addAndroidIconLinks()
	hdw.addFaviconLinks()
	hdw.addAppleIconLinks()
	hdw.addStandardMeta()
	hdw.addGoogleApiLinkToJQuery()
}

func (hdw *htmlDocWrapper) addGoogleApiLinkToJQuery() {
	hdw.htmlDoc.AddToHead(createNode("script").Attr("src", "https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"))
}

func (hdw *htmlDocWrapper) addTitle(txt string) {
	hdw.htmlDoc.AddToHead(createNode("title").AppendText(txt))
}

func (hdw *htmlDocWrapper) addStandardMeta() {
	name_metas := []string{
		"viewport", "width=device-width, initial-scale=1.0",
		"robots", "index,follow",
		"author", "Ingmar Drewing",
		"publisher", "Ingmar Drewing",
		"keywords", "web comic, comic, cartoon, sci fi, satire, parody, science fiction, action, software industry, pulp, nerd, geek",
		"DC.Subject", "web comic, comic, cartoon, sci fi, science fiction, satire, parody action, software industry",
		"page-topic", "Science Fiction Web-Comic",
	}
	hdw.addNameValueMetas(name_metas)
	hdw.htmlDoc.AddToHead(createNode("meta").Attr("http-equiv", "content-type").Attr("content", "text/html;charset=UTF-8"))
}

func (hdw *htmlDocWrapper) addNameValueMetas(metaData []string) {
	for i := 0; i < len(metaData); i += 2 {
		m := createNode("meta")
		m.Attr(metaData[i], metaData[i+1])
		hdw.htmlDoc.AddToHead(m)
	}
}

func (hdw *htmlDocWrapper) addCopyrightNotifier(year string) {
	hdw.htmlDoc.AddToBody(createNode("div").Attr("class", "copyright").AppendText(`All content including but not limited to the art, characters, story, website design & graphics are &copy; copyright 2013-` + year + ` Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.`))
}

func (hdw *htmlDocWrapper) addCookieLawInfo() {
	hdw.htmlDoc.AddToBody(createNode("div").Attr("id", "cookie-law-info-bar").AppendText(`This website uses cookies to improve your experience. We'll assume you're ok with this, but you can opt-out if you wish.<a href="#" id="cookie_action_close_header" class="medium cli-plugin-button cli-plugin-main-button">Accept</a> <a href="http://www.drewing.de/blog/impressum-imprint/" id="CONSTANT_OPEN_URL" target="_blank" class="cli-plugin-main-link">Read More</a>`))
}

func (hdw *htmlDocWrapper) addFooterNavi(navi string) {
	n := createNode("footer")
	n.AppendTag("nav").AppendText(navi)
	hdw.htmlDoc.AddToBody(n)
}

func (hdw *htmlDocWrapper) addFaviconLinks() {
	iconSizes := []string{
		"32x32",
		"96x96",
		"16x16",
	}
	for _, s := range iconSizes {
		l := createNode("link")
		l.Attr("rel", "icon")
		l.Attr("type", "image/png")
		l.Attr("sizes", s)
		l.Attr("href", "/favicon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) addAndroidIconLinks() {
	androidSizes := []string{
		"192x192",
	}
	for _, s := range androidSizes {
		l := createNode("link")
		l.Attr("rel", "icon")
		l.Attr("type", "image/png")
		l.Attr("sizes", s)
		l.Attr("href", "/android-icon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) addAppleIconLinks() {
	appleSizes := []string{
		"57x57",
		"60x60",
		"72x72",
		"76x76",
		"114x114",
		"120x120",
		"144x144",
		"152x152",
		"180x180",
	}
	for _, s := range appleSizes {
		l := createNode("link")
		l.Attr("rel", "apple-touch-icon")
		l.Attr("sizes", s)
		l.Attr("href", "/apple-icon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) Render() string {
	return hdw.htmlDoc.Render()
}
