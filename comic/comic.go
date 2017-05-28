package comic

import "github.com/ingmardrewing/gomic/config"

type Comic struct {
	rootpath string
	pages    []*Page
}

func NewComic(pages []*Page) Comic {
	c := Comic{config.Rootpath(), pages}
	return c
}

func (c *Comic) AddPage(p *Page) {
	c.pages = append(c.pages, p)
}

func (c *Comic) Get10LastComicPagesNewestFirst() []*Page {
	// get splice with last 10 pages
	last10 := c.pages[len(c.pages)-11:]

	// reorder them, so last issued is first in splice
	for i := len(last10)/2 - 1; i >= 0; i-- {
		opp := len(last10) - 1 - i
		last10[i], last10[opp] = last10[opp], last10[i]
	}
	return last10
}

func (c *Comic) ConnectPages() {
	for i, p := range c.pages {
		p.SetRels(
			c.firstFor(i),
			c.previousFor(i),
			c.nextFor(i),
			c.lastFor(i))
	}
}

func (c *Comic) GetPages() []*Page {
	return c.pages
}

func (c *Comic) pageIndexExists(i int) bool {
	l := len(c.pages)
	return l > 0 && i >= 0 && i < l
}

func (c *Comic) firstFor(i int) *Page {
	if c.pageIndexExists(0) && i != 0 {
		return c.firstPage()
	}
	return nil
}

func (c *Comic) previousFor(i int) *Page {
	if c.pageIndexExists(i - 1) {
		return c.pages[i-1]
	}
	return nil
}

func (c *Comic) nextFor(i int) *Page {
	if c.pageIndexExists(i + 1) {
		return c.pages[i+1]
	}
	return nil
}

func (c *Comic) lastFor(i int) *Page {

	l := len(c.pages)
	if l > 0 && i != l-1 {
		return c.LastPage()
	}
	return nil
}

func (c *Comic) firstPage() *Page {
	if c.pageIndexExists(0) {
		return c.pages[0]
	}
	return nil
}

func (c *Comic) LastPage() *Page {
	l := len(c.pages)
	if l > 0 {
		return c.pages[l-1]
	}
	return nil
}
