package page

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ingmardrewing/gomic/config"
)

type Page struct {
	title, path, imgUrl, disqusId, act string
	first, prev, next, last            *Page
	meta, navi                         [][]string
}

func NewPageFromFilename(filename string) *Page {

	var title, path, imgUrl, disqusId, act string
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter title for %s: ", filename)
		title, _ = reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Printf("Enter act for %s: ", filename)
		act, _ = reader.ReadString('\n')
		act = strings.TrimSpace(act)

		whitespace := regexp.MustCompile(`\s+`)
		forbidden := regexp.MustCompile(`[^-A-Za-z0-9]`)
		trailingdash := regexp.MustCompile(`-$`)
		pathTitle := whitespace.ReplaceAllString(title, "-")
		pathTitle = forbidden.ReplaceAllString(pathTitle, "")
		pathTitle = trailingdash.ReplaceAllString(pathTitle, "")

		t := time.Now()
		y := t.Year()
		m := int(t.Month())
		d := t.Day()
		path = fmt.Sprintf("/%d/%02d/%02d/%s", y, m, d, pathTitle)

		id := y*10000 + m*100 + d
		disqusId = fmt.Sprintf("%d https://DevAbo.de/?p=%d", id, id)

		imgUrl = fmt.Sprintf("https://s3-us-west-1.amazonaws.com/devabode-us/comicstrips/%s", filename)

		summary := fmt.Sprintf("\ntitle: %s\npath: %s\ndisqusId: %s\nimgUrl: %s\n", title, path, disqusId, imgUrl)

		answer := AskUser(
			fmt.Sprintf(
				"Creating the following page:\n%s\nok? [yN]", summary))

		if answer {
			break
		}
		fmt.Println("Okay, let's try again ...")
	}

	return &Page{title, path, imgUrl, disqusId, act, nil, nil, nil, nil, [][]string{}, [][]string{}}
}

func AskUser(question string) bool {
	fmt.Println(question)
	reader := bufio.NewReader(os.Stdin)
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)
	return confirmation == "y" || confirmation == "Y"
}

func NewPage(
	title string,
	path string,
	imgUrl string,
	disqusId string,
	act string) *Page {
	return &Page{title, path, imgUrl, disqusId, act,
		nil, nil, nil, nil, [][]string{}, [][]string{}}

}

func (p *Page) ImageFilename() string {
	pathParts := strings.Split(p.imgUrl, "/")
	return pathParts[len(pathParts)-1]
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) DisqusId() string {
	return p.disqusId
}

func (p *Page) ImgUrl() string {
	return p.imgUrl
}

func (p *Page) ThumnailUrl() string {
	thumbUrl := fmt.Sprintf("https://s3-us-west-1.amazonaws.com/devabode-us/%s/thumb_%s", config.AwsDir(), p.ImageFilename())
	return thumbUrl
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

func (p *Page) Date() string {
	parts := strings.Split(p.FSPath(), "/")
	loc, _ := time.LoadLocation("Europe/Berlin")
	y, _ := strconv.Atoi(parts[1])
	m, _ := strconv.Atoi(parts[2])
	d, _ := strconv.Atoi(parts[3])
	date := time.Date(y, time.Month(m), d, 20, 0, 0, 0, loc)
	return date.Format(time.RFC1123Z)
}

func (p *Page) Act() string {
	return p.act
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
