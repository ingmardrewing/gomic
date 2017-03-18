package page

import "fmt"

type Page struct {
	title, url              string
	first, prev, next, last *Page
}

func NewPage(title string, url string) *Page {
	return &Page{title, url, nil, nil, nil, nil}
}

func (p *Page) SetRels(first *Page, prev *Page, next *Page, last *Page) {
	p.first = first
	p.prev = prev
	p.next = next
	p.last = last
}

func (p *Page) Html() string {
	return fmt.Sprintf(htmlFormat, p.title, p.meta(), p.img(), p.navi())
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
		return fmt.Sprintf(headerLinkFormat, rel, linked.title, linked.url)
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
		return fmt.Sprintf(navLinkFormat, rel, linked.title, linked.url, label)
	}
	return ""
}

func (p *Page) img() string {
	if p.next != nil {
		img := fmt.Sprintf(imageFormat, p.url)
		return fmt.Sprintf(imageWrapperFormat, p.next.url, p.next.title, img)
	}
	return fmt.Sprintf(imageFormat, p.url)
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
	</head>
	<body>
		%s
		%s
	</body>
</html>
`
