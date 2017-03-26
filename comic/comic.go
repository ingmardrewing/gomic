package comic

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/ingmardrewing/gomic/aws"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/db"
	"github.com/ingmardrewing/gomic/page"
)

type Comic struct {
	rootpath string
	pages    []*page.Page
}

func (c *Comic) generatePages(rows *sql.Rows) {
	for rows.Next() {
		var (
			title    sql.NullString
			path     sql.NullString
			imgUrl   sql.NullString
			disqusId sql.NullString
			act      sql.NullString
			id       sql.NullInt64
		)
		rows.Scan(&title, &path, &imgUrl, &disqusId, &act, &id)
		p := page.NewPage(title.String, path.String, imgUrl.String, disqusId.String, act.String)
		c.AddPage(p)
	}
}

func NewComic(rows *sql.Rows) Comic {
	pages := []*page.Page{}
	c := Comic{config.Rootpath(), pages}
	c.generatePages(rows)
	return c
}

func (c *Comic) AddPage(p *page.Page) {
	c.pages = append(c.pages, p)
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
		return c.LastPage()
	}
	return nil
}

func (c *Comic) firstPage() *page.Page {
	if c.pageIndexExists(0) {
		return c.pages[0]
	}
	return nil
}

func (c *Comic) LastPage() *page.Page {
	l := len(c.pages)
	if l > 0 {
		return c.pages[l-1]
	}
	return nil
}

func (c *Comic) isRelevant(filename string) bool {
	irr := ".DS_Store"
	if filename == irr {
		return false
	}
	thumb := regexp.MustCompile(`^thumb_`)
	if thumb.MatchString(filename) {
		return false
	}
	return true
}

func (c *Comic) isNewFile(filename string) bool {
	if !c.isRelevant(filename) {
		return false
	}
	for _, p := range c.pages {
		fn := p.ImageFilename()
		if fn == filename {
			return false
		}
	}
	return true
}

func (c *Comic) CheckForNewPages(filenames []string) {
	for _, f := range filenames {
		if c.isNewFile(f) {
			log.Printf("Found new file: %s", f)
			p := page.NewPageFromFilename(f)
			aws.UploadPage(p)
			db.InsertPage(p)
			c.AddPage(p)
		}
	}
}
