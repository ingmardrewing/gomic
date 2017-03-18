package page

import (
	"fmt"
	"time"
)

type Page struct {
	title, path, imgUrl, servedrootpath string
	first, prev, next, last             *Page
}

func NewPage(
	title string,
	path string,
	imgUrl string,
	servedrootpath string) *Page {
	return &Page{title, path, imgUrl, servedrootpath, nil, nil, nil, nil}
}

func (p *Page) SetRels(first *Page, prev *Page, next *Page, last *Page) {
	p.first = first
	p.prev = prev
	p.next = next
	p.last = last
}

func (p *Page) version() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02dT%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (p *Page) Html() string {
	return fmt.Sprintf(
		htmlFormat, p.title, p.meta(),
		p.version(), p.img(), p.navi())
}

func (p *Page) meta() string {
	meta := ""
	meta += p.getHeaderLink("start", p.first)
	meta += p.getHeaderLink("prev", p.prev)
	meta += p.getHeaderLink("next", p.next)
	meta += p.getHeaderLink("last", p.last)
	return meta
}

func (p *Page) getHeaderLink(rel string, linked *Page) string {
	if linked != nil {
		return fmt.Sprintf(headerLinkFormat, rel, linked.title, linked.Path())
	}
	return ""
}

func (p *Page) navi() string {
	n := ""
	n += p.getNavLink("first", "&lt;&lt; first", p.first)
	n += p.getNavLink("previous", "&lt; previous", p.prev)
	n += p.getNavLink("next", "&gt; next", p.next)
	n += p.getNavLink("last", "&gt;&gt; newest", p.last)
	return fmt.Sprintf(navWrapperFormat, n)
}

func (p *Page) getNavLink(rel string, label string, linked *Page) string {
	if linked != nil {
		return fmt.Sprintf(navLinkFormat, rel, linked.title, linked.Path(), label)
	}
	return ""
}

func (p *Page) Path() string {
	path := p.servedrootpath + p.path
	return path
}

func (p *Page) FSPath() string {
	return p.path
}

func (p *Page) img() string {
	if p.next != nil {
		img := fmt.Sprintf(imageFormat, p.imgUrl)
		return fmt.Sprintf(imageWrapperFormat, p.next.Path(), p.next.title, img)
	}
	return fmt.Sprintf(imageFormat, p.imgUrl)
}

const imageWrapperFormat = `<a href="%s" rel="next" title="%s">%s</a>`
const imageFormat = `<img src="%s" width="800" height="1334" alt="">`
const headerLinkFormat = `<link rel="%s" title="%s" href="%s">`
const navLinkFormat = `<a rel="%s" title="%s" href="%s">%s</a>`
const navWrapperFormat = `<nav>%s</nav>`
const htmlFormat = `<!doctype html>
<html>
	<head>
		<title>%s</title>
		%s
		<link rel="stylesheet" href="/~drewing/gomic/css/style.css?version=%s" type="text/css">
	</head>
	<body>
		%s
		%s
	</body>
</html>
`
