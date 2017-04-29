package comic

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/ingmardrewing/gomic/config"
)

type Comic struct {
	rootpath string
	pages    []*Page
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
		p := NewPage(title.String, path.String, imgUrl.String, disqusId.String, act.String)
		c.AddPage(p)
	}
}

func NewComic(rows *sql.Rows) Comic {
	pages := []*Page{}
	c := Comic{config.Rootpath(), pages}
	c.generatePages(rows)
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

func (c *Comic) isRelevant(filename string) bool {
	irr := ".DS_Store"
	if filename == irr {
		fmt.Println(1)
		return false
	}
	thumb := regexp.MustCompile(`^thumb_`)
	if thumb.MatchString(filename) {
		fmt.Println(2)
		return false
	}
	fmt.Println(3)
	return true
}

func (c *Comic) IsNewFile(filename string) bool {
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

func (c *Comic) CreateThumbnail(filename string) {
	command := "/Users/drewing/bin/createDevabodeThumbFromPath.pl"
	thumbnail_px_width := "150"
	args := []string{config.PngDir() + filename, thumbnail_px_width}
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Created Thumbnail for %s\n", filename)
}
