package comic

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ingmardrewing/gomic/config"
)

type Pages struct {
	Pages []*Page
}

type Page struct {
	Id, PageNumber                                  int
	Title, Description, Path, ImgUrl, DisqusId, Act string
	first, prev, next, last                         *Page
	meta, navi                                      [][]string
}

func getUserInput(prompt string) string {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func createPathTitleFromTitle(title string) string {
	whitespace := regexp.MustCompile(`\s+`)
	forbidden := regexp.MustCompile(`[^-A-Za-z0-9]`)
	trailingdash := regexp.MustCompile(`-$`)

	pathTitle := whitespace.ReplaceAllString(title, "-")
	pathTitle = forbidden.ReplaceAllString(pathTitle, "")
	return trailingdash.ReplaceAllString(pathTitle, "")
}

func getYMD() (int, int, int) {
	t := time.Now()
	return t.Year(), int(t.Month()), t.Day()
}

func getPath(title string, y int, m int, d int) string {
	pathTitle := createPathTitleFromTitle(title)
	return fmt.Sprintf("/%d/%02d/%02d/%s", y, m, d, pathTitle)
}

func getDisqusId(y int, m int, d int) string {
	id := y*10000 + m*100 + d
	disqusId := fmt.Sprintf("%d https://DevAbo.de/?p=%d", id, id)
	return disqusId
}

func getPageData(filename string) (string, string, string, string, string, string) {
	act := getUserInput("Enter act for " + filename + ": ")
	title := getUserInput("Enter title for " + filename + ": ")
	description := getUserInput("Enter description for " + filename + ": ")
	y, m, d := getYMD()
	path := getPath(title, y, m, d)
	disqusId := getDisqusId(y, m, d)
	imgUrl := fmt.Sprintf("https://s3-us-west-1.amazonaws.com/devabode-us/comicstrips/%s", filename)
	return act, title, path, disqusId, imgUrl, description
}

func getPageFromFilenameAndUserInput(filename string) *Page {
	act, title, path, disqusId, imgUrl, description := getPageData(filename)
	return &Page{0, 0, title, description, path, imgUrl, disqusId, act, nil, nil, nil, nil, [][]string{}, [][]string{}}
}

func CreateThumbnail(filename string) {
	command := "/Users/drewing/bin/createDevabodeThumbFromPath.pl"
	pngDir := config.PngDir() + filename
	thumbnail_px_width := "150"
	args := []string{pngDir, thumbnail_px_width}
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Created Thumbnail for %s\n", filename)
}

func NewPageFromFilename(filename string) *Page {
	for {
		page := getPageFromFilenameAndUserInput(filename)
		summary := fmt.Sprintf("\ntitle: %s\ndescription: %s\npath: %s\ndisqusId: %s\nimgUrl: %s\n", page.Title, page.Description, page.Path, page.DisqusId, page.ImgUrl)
		answer := AskUser(fmt.Sprintf("Creating the following page:\n%s\nok? [yN]", summary))
		if answer {
			return page
		}
		fmt.Println("Okay, let's try again ...")
	}
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
	description string,
	path string,
	imgUrl string,
	disqusId string,
	act string) *Page {
	return &Page{0, 0, title, description, path, imgUrl, disqusId, act,
		nil, nil, nil, nil, [][]string{}, [][]string{}}
}

func (p *Page) GetImageFilename() string {
	pathParts := strings.Split(p.ImgUrl, "/")
	return pathParts[len(pathParts)-1]
}

func (p *Page) GetProdUrl() string {
	return "https://devabo.de" + p.Path
}

func (p *Page) GetTitle() string {
	return p.Title
}

func (p *Page) GetDisqusId() string {
	return p.DisqusId
}

func (p *Page) GetImgUrl() string {
	return p.ImgUrl
}

func (p *Page) GetDescription() string {
	return p.Description
}

func (p *Page) GetThumnailUrl() string {
	thumbUrl := fmt.Sprintf("https://s3-us-west-1.amazonaws.com/devabode-us/%s/thumb_%s", config.AwsDir(), p.GetImageFilename())
	return thumbUrl
}

func (p *Page) GetDisqusIdentifier() string {
	return p.DisqusId
}

func (p *Page) SetRels(first *Page, prev *Page, next *Page, last *Page) {
	p.first = first
	p.prev = prev
	p.next = next
	p.last = last
}

func (p *Page) fillMeta() {
	if p.first != nil {
		p.addMeta("start", p.first.Title, p.first.GetPath())
	}
	if p.prev != nil {
		p.addMeta("prev", p.prev.Title, p.prev.GetPath())
	}
	if p.next != nil {
		p.addMeta("next", p.next.Title, p.next.GetPath())
	}
	if p.last != nil {
		p.addMeta("last", p.last.Title, p.last.GetPath())
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
		p.addNavi("first", p.first.Title, p.first.GetPath(), "&lt;&lt; first")
	}
	if p.prev != nil {
		p.addNavi("previous", p.prev.Title, p.prev.GetPath(), "&lt; previous")
	}
	if p.next != nil {
		p.addNavi("next", p.next.Title, p.next.GetPath(), "next &gt;")
	}
	if p.last != nil {
		p.addNavi("last", p.last.Title, p.last.GetPath(), "newest &gt;")
	}
}

func (p *Page) GetDateFromFSPath() string {
	parts := strings.Split(p.FSPath(), "/")
	loc, _ := time.LoadLocation("Europe/Berlin")
	y, _ := strconv.Atoi(parts[1])
	m, _ := strconv.Atoi(parts[2])
	d, _ := strconv.Atoi(parts[3])
	date := time.Date(y, time.Month(m), d, 20, 0, 0, 0, loc)
	return date.Format(time.RFC1123Z)
}

func (p *Page) GetAct() string {
	return p.Act
}

func (p *Page) IsLast() bool {
	return p.last == nil
}

func (p *Page) UrlToNext() string {
	return p.next.GetPath()
}

func (p *Page) addNavi(rel string, label string, title string, path string) {
	n := []string{rel, label, title, path}
	p.navi = append(p.navi, n)
}

func (p *Page) GetNavi() [][]string {
	p.fillNavi()
	return p.navi
}

func (p *Page) GetPath() string {
	path := config.Servedrootpath() + p.Path
	return path
}

func (p *Page) FSPath() string {
	return p.Path
}
