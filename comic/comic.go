package comic

import (
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
)

type Comic struct {
	rootpath string
	pages    []*page.Page
}

func (c *Comic) generatePages() {
	for _, pc := range config.Pages() {
		p := c.generatePage(pc, config.Servedrootpath())
		c.AddPage(p)
	}
}

func (c *Comic) generatePage(pc map[string]string, servedRootPath string) *page.Page {
	return page.NewPage(pc["title"], pc["path"], pc["imgUrl"], pc["disqusId"], servedRootPath)
}

func NewComic() Comic {
	pages := []*page.Page{}
	c := Comic{config.Rootpath(), pages}
	c.generatePages()
	c.connectPages()
	return c
}

func (c *Comic) AddPage(p *page.Page) {
	c.pages = append(c.pages, p)
}

func (c *Comic) connectPages() {
	for i, p := range c.pages {
		p.SetRels(
			c.firstFor(i),
			c.previousFor(i),
			c.nextFor(i),
			c.lastFor(i))
	}
}

func (c *Comic) GetPages() []*page.Page {
	return c.pages
}

func (c *Comic) pageIndexExists(i int) bool {
	l := len(c.pages)
	return l > 0 && i >= 0 && i < l
}

func (c *Comic) firstFor(i int) *page.Page {
	if c.pageIndexExists(0) && i != 0 {
		return c.firstPage()
	}
	return nil
}

func (c *Comic) previousFor(i int) *page.Page {
	if c.pageIndexExists(i - 1) {
		return c.pages[i-1]
	}
	return nil
}

func (c *Comic) nextFor(i int) *page.Page {
	if c.pageIndexExists(i + 1) {
		return c.pages[i+1]
	}
	return nil
}

func (c *Comic) lastFor(i int) *page.Page {
	l := len(c.pages)
	if l > 0 && i != l-1 {
		return c.lastPage()
	}
	return nil
}

func (c *Comic) firstPage() *page.Page {
	if c.pageIndexExists(0) {
		return c.pages[0]
	}
	return nil
}

func (c *Comic) lastPage() *page.Page {
	l := len(c.pages)
	if l > 0 {
		return c.pages[l-1]
	}
	return nil
}
