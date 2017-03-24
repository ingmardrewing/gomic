package page

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ingmardrewing/gomic/config"
)

type Page struct {
	title, path, imgUrl, disqusId string
	first, prev, next, last       *Page
	meta, navi                    [][]string
}

func NewPageFromFilename(filename string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter title for %s: ", filename)
	title, _ := reader.ReadString('\n')

	whitespace := regexp.MustCompile(`\s+`)
	forbidden := regexp.MustCompile(`[^-A-Za-z0-9]`)
	pathTitle := whitespace.ReplaceAllString(title, "-")
	pathTitle = forbidden.ReplaceAllString(pathTitle, "")

	t := time.Now()
	y := t.Year()
	m := int(t.Month())
	d := t.Day()
	path := fmt.Sprintf("/%d/%02d/%02d/%s", y, m, d, pathTitle)
	log.Println(path)
}

func NewPage(
	title string,
	path string,
	imgUrl string,
	disqusId string) *Page {
	return &Page{title, path, imgUrl, disqusId,
		nil, nil, nil, nil, [][]string{}, [][]string{}}

}

func (p *Page) Filename() string {
	pathParts := strings.Split(p.imgUrl, "/")
	return pathParts[len(pathParts)-1]
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) DisqusIdentifier() string {
	return p.disqusId
}

func (p *Page) SetRels(first *Page, prev *Page, next *Page, last *Page) {
	p.first = first
	p.prev = prev
	p.next = next
	p.last = last
}

func (p *Page) fillMeta() {
	if p.first != nil {
		p.addMeta("start", p.first.title, p.first.Path())
	}
	if p.prev != nil {
		p.addMeta("prev", p.prev.title, p.prev.Path())
	}
	if p.next != nil {
		p.addMeta("next", p.next.title, p.next.Path())
	}
	if p.last != nil {
		p.addMeta("last", p.last.title, p.last.Path())
	}
}

func (p *Page) addMeta(rel string, title string, path string) {
	l := []string{rel, title, path}
	p.meta = append(p.meta, l)
}

func (p *Page) GetMeta() [][]string {
	p.fillMeta()
	return p.meta
}

func (p *Page) fillNavi() {
	if p.first != nil {
		p.addNavi("first", p.first.title, p.first.Path(), "&lt;&lt; first")
	}
	if p.prev != nil {
		p.addNavi("previous", p.prev.title, p.prev.Path(), "&lt; previous")
	}
	if p.next != nil {
		p.addNavi("next", p.next.title, p.next.Path(), "next &gt;")
	}
	if p.last != nil {
		p.addNavi("last", p.last.title, p.last.Path(), "newest &gt;")
	}
}

func (p *Page) IsLast() bool {
	return p.last == nil
}

func (p *Page) UrlToNext() string {
	return p.next.Path()
}

func (p *Page) addNavi(rel string, label string, title string, path string) {
	n := []string{rel, label, title, path}
	p.navi = append(p.navi, n)
}

func (p *Page) GetNavi() [][]string {
	p.fillNavi()
	return p.navi
}

func (p *Page) Path() string {
	path := config.Servedrootpath() + p.path
	return path
}

func (p *Page) FSPath() string {
	return p.path
}

func (p *Page) Img() string {
	return p.imgUrl
}
