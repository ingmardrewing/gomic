package fs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
	"github.com/nfnt/resize"
)

func main() {
	fmt.Println("vim-go")
}

func ReadImageFilenames() []string {
	path := config.PngDir()
	files, _ := ioutil.ReadDir(path)
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames
}

type Output struct {
	comic *comic.Comic
}

func NewOutput(comic *comic.Comic) *Output {
	return &Output{comic}
}

func (o *Output) WriteToFilesystem() {
	o.writeNarrativePages()
	o.writeCss()
	o.writeJs()
	o.writeArchive()
	o.writeRss()
	o.writeAbout()
	o.writeImprint()
}

func (o *Output) writeAbout() {
	ah := NewDataHtml(about, config.Servedrootpath()+"/about.html")
	o.writeStringToFS(config.Rootpath()+"/about.html", ah.writePage("About"))
}

func (o *Output) writeImprint() {
	ah := NewDataHtml(imprint, config.Servedrootpath()+"/imprint.html")
	o.writeStringToFS(config.Rootpath()+"/imprint.html", ah.writePage("Imprint"))
}

func (o *Output) writeRss() {
	rss := newRss(o.comic)
	path := config.Rootpath() + "/feed/rss.xml"
	log.Println("Writing rss: ", path)
	o.writeStringToFS(path, rss.Rss())
}

func (o *Output) writeNarrativePages() {
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
}

func (o *Output) writeThumbnailFor(p *page.Page) string {
	imgpath := config.PngDir() + p.ImageFilename()
	outimgpath := config.PngDir() + "thumb_" + p.ImageFilename()
	if _, err := os.Stat(outimgpath); os.IsNotExist(err) {
		// open "test.jpg"
		file, err := os.Open(imgpath)
		if err != nil {
			log.Fatal(err)
		}

		// decode jpeg into image.Image
		img, err := png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(150, 0, img, resize.Lanczos3)

		out, err := os.Create(outimgpath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		png.Encode(out, m)
	}
	return outimgpath
}

func (o *Output) getBase64FromPngFile(path string) (string, int, int) {
	imgFile, err := os.Open(path) // a QR code image

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	b := base64.StdEncoding.EncodeToString(buf)

	imgFile2, err := os.Open(path) // a QR code image

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile2.Close()
	ime, _, err := image.DecodeConfig(imgFile2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
	}

	return b, ime.Width, ime.Height
}

func (o *Output) writeArchive() {
	list := []string{}
	for _, p := range o.comic.GetPages() {
		path := o.writeThumbnailFor(p)
		b, w, h := o.getBase64FromPngFile(path)
		list = append(list, fmt.Sprintf(`<li><a href="%s"><img src="data:image/png;base64,%s" width="%d" height="%d" alt="%s" title="%s"></a></li>`, p.Path(), b, w, h, p.Title(), p.Title()))

	}

	arc := fmt.Sprintf(`<ul class="archive">%s</ul>`, strings.Join(list, "\n"))
	ah := NewDataHtml(arc, config.Servedrootpath()+"/archive.html")
	o.writeStringToFS(config.Rootpath()+"/archive.html", ah.writePage("Archive"))
}

func (o *Output) writeCss() {
	p := config.Rootpath() + "/css"
	o.prepareFileSystem(p)
	fp := p + "/style.css"
	o.writeStringToFS(fp, css)
}

func (o *Output) writeJs() {
	p := config.Rootpath() + "/js"
	o.prepareFileSystem(p)
	fp := p + "/script.js"
	o.writeStringToFS(fp, js)
}

func (o *Output) writePageToFileSystem(p *page.Page) {
	absPath := config.Rootpath() + p.FSPath()
	o.prepareFileSystem(absPath)

	h := NewNarrativePageHtml(p)
	html := h.writePage()
	o.writeStringToFS(absPath+"/index.html", html)
	if p.IsLast() {
		o.writeStringToFS(config.Rootpath()+"/index.html", html)
	}
}

func (o *Output) writeStringToFS(absPath string, html string) {
	//log.Println("writing html to filesystem: ", absPath)
	b := []byte(html)
	err := ioutil.WriteFile(absPath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func (o *Output) prepareFileSystem(absPath string) {
	exists, err := o.pathExists(absPath)
	if err != nil {
		panic(err.Error())
	}
	if !exists {
		log.Println("creating path", absPath)
		os.MkdirAll(absPath, 0755)
	}
}

func (o *Output) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type HTML struct{}

func (html *HTML) version() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02dT%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (html *HTML) getFooterNavi() string {
	p := config.Servedrootpath()
	h := ""
	h += fmt.Sprintf(`<a href="http://twitter.com/devabo_de">Twitter</a>
	<a href="%s/about.html">About</a>
	<a href="%s/feed/rss.xml">RSS</a>
	<a href="%s/archive.html">Archive</a>
	<a href="%s/imprint.html">Imprint / Impressum</a>
	`, p, p, p, p)
	if config.IsProd() {
		h += analytics
	}
	return h
}

func (html *HTML) getCssLink() string {
	path := config.Servedrootpath() + "/css/style.css?version=" + html.version()
	format := `<link rel="stylesheet" href="%s" type="text/css">`
	return fmt.Sprintf(format, path)
}

func (html *HTML) getJsLink() string {
	path := config.Servedrootpath() + "/js/script.js?version=" + html.version()
	format := `<script src="%s" type="text/javascript" language="javascript"></script>`
	return fmt.Sprintf(format, path)
}

func (html *HTML) getHeaderLink(vals ...string) string {
	l := fmt.Sprintf(`<link rel="%s" title="%s" href="%s">`, vals[0], vals[1], vals[2])
	return l
}

func (html *HTML) getHeadline(txt string) string {
	return fmt.Sprintf(`<h3>%s</h3>`, txt)
}

func (html *HTML) getMetaHtml() string {
	return ""
}

func (html *HTML) getNaviHtml() string {
	return ""
}

func (html *HTML) getTitle() string {
	return ""
}

func (html *HTML) getHeaderHtml() string {
	hl := html.getHeadline("")
	s := config.Servedrootpath()
	return fmt.Sprintf(`
	<a href="%s" class="home"><!--DevAbo.de--></a>
    <a href="%s/2013/08/01/a-step-in-the-dark/" class="orange">New Reader? Start here!</a>
	%s`, s, s, hl)
}

func (html *HTML) getContent() string {
	return "x"
}

func (html *HTML) writePage() string {
	hdw := newHtmlDocWrapper()
	hdw.Init()

	css_path := config.Servedrootpath() + "/css/style.css?version=" + hdw.Version()
	hdw.AddToHead(createNode("link").Attr("rel", "stylesheet").Attr("href", css_path).Attr("type", "text/css"))

	js_path := config.Servedrootpath() + "/js/script.js?version=" + hdw.Version()
	hdw.AddToHead(createNode("script").Attr("src", js_path).Attr("type", "text/javascript").Attr("language", "javascript"))
	hdw.AddTitle("DevAbo.de | Graphic Novel")

	header := createNode("header").AppendText(html.getHeaderHtml())
	hdw.AddToBody(header)
	hdw.AddToBody(createText("<!-- have you not -->"))

	main := createNode("main")
	main.AppendText(html.getContent())
	main.AppendText(html.getNaviHtml())
	hdw.AddToBody(main)

	hdw.AddCopyrightNotifier(strconv.Itoa(time.Now().Year()))

	hdw.AddFooterNavi(html.getFooterNavi())

	return hdw.Render()
}

type DataHtml struct {
	HTML
	content string
	url     string
}

func NewDataHtml(content string, url string) *DataHtml {
	return &DataHtml{HTML{}, content, url}
}

func (ah *DataHtml) getContent() string {
	return ah.content
}

func (ah *DataHtml) writePage(title string) string {
	hdw := newHtmlDocWrapper()
	hdw.Init()

	css_path := config.Servedrootpath() + "/css/style.css?version=" + hdw.Version()
	hdw.AddToHead(createNode("link").Attr("rel", "stylesheet").Attr("href", css_path).Attr("type", "text/css"))

	js_path := config.Servedrootpath() + "/js/script.js?version=" + hdw.Version()
	hdw.AddToHead(createNode("script").Attr("src", js_path).Attr("type", "text/javascript").Attr("language", "javascript"))

	hdw.AddTitle("DevAbo.de | Graphic Novel | " + title)

	header := createNode("header").AppendText(ah.getHeaderHtml())
	hdw.AddToBody(header)

	main := createNode("main")
	main.AppendText(ah.getContent())
	hdw.AddToBody(main)

	hdw.AddCopyrightNotifier(strconv.Itoa(time.Now().Year()))

	hdw.AddFooterNavi(ah.getFooterNavi())

	return hdw.Render()
}

type NarrativePageHtml struct {
	HTML
	p          *page.Page
	title      string
	meta       string
	csslink    string
	img        string
	navi       string
	footerNavi string
}

func NewNarrativePageHtml(p *page.Page) *NarrativePageHtml {
	return &NarrativePageHtml{HTML{}, p, "", "", "", "", "", ""}
}

func (h *NarrativePageHtml) writePage() string {
	hdw := newHtmlDocWrapper()
	hdw.Init()

	css_path := config.Servedrootpath() + "/css/style.css?version=" + hdw.Version()
	hdw.AddToHead(createNode("link").Attr("rel", "stylesheet").Attr("href", css_path).Attr("type", "text/css"))
	hdw.AddToHead(createNode("link").Attr("rel", "canonical").Attr("href", h.p.Path()))

	js_path := config.Servedrootpath() + "/js/script.js?version=" + hdw.Version()
	hdw.AddToHead(createNode("script").Attr("src", js_path).Attr("type", "text/javascript").Attr("language", "javascript"))
	hdw.AddTitle("DevAbo.de | Graphic Novel | " + h.p.Title())

	header := createNode("header").AppendText(h.getHeaderHtml())
	hdw.AddToBody(header)

	main := createNode("main")
	main.AppendText(h.getContent())
	main.AppendText(h.getNaviHtml())
	main.AppendText(h.getDisqus())
	hdw.AddToBody(main)

	hdw.AddCopyrightNotifier(strconv.Itoa(time.Now().Year()))

	hdw.AddFooterNavi(h.getFooterNavi())

	hdw.AddToHead(createNode("meta").Attr("property", "og:title").Attr("content", h.p.Title()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:url").Attr("content", h.p.Path()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:image").Attr("content", h.p.ImgUrl()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:description").Attr("content", "A dystopian science-fiction webcomic set 1337 years after WW III"))
	hdw.AddToHead(createNode("meta").Attr("property", "og:site_name").Attr("content", "DevAbo.de"))
	hdw.AddToHead(createNode("meta").Attr("property", "og:type").Attr("content", "article"))
	hdw.AddToHead(createNode("meta").Attr("property", "article:published_time").Attr("content", h.p.Date()))
	hdw.AddToHead(createNode("meta").Attr("property", "article:modified_time").Attr("content", h.p.Date()))
	hdw.AddToHead(createNode("meta").Attr("property", "article:section").Attr("content", "Science-Fiction"))
	hdw.AddToHead(createNode("meta").Attr("property", "article:tag").Attr("content", "comic, graphic novel, webcomic, science-fiction, sci-fi"))

	hdw.AddToHead(createNode("meta").Attr("itemprop", "name").Attr("content", h.p.Title()))
	hdw.AddToHead(createNode("meta").Attr("itemprop", "name").Attr("description", "A dystopian sci-fi webcomic about the life of software developers"))
	hdw.AddToHead(createNode("meta").Attr("itemprop", "image").Attr("content", h.p.ImgUrl()))

	hdw.AddToHead(createNode("meta").Attr("name", "twitter:card").Attr("content", h.p.ImgUrl()))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:site").Attr("content", "@devabo_de"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:title").Attr("content", h.p.Title()))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:description").Attr("content", "A dystopian science-fiction webcomic set 1337 years after WW III"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:creator").Attr("content", "@ingmardrewing"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:image:src").Attr("content", h.p.ImgUrl()))

	return hdw.Render()
}

func (h *NarrativePageHtml) getContent() string {
	f := `<img src="%s" width="800" height="1334" alt="">`
	html := fmt.Sprintf(f, h.p.Img())
	if !h.p.IsLast() {
		html = fmt.Sprintf(`<a href="%s">%s</a>`, h.p.UrlToNext(), html)
	}
	return html
}

func (h *NarrativePageHtml) getNaviHtml() string {
	ns := h.p.GetNavi()
	html := ""
	for _, n := range ns {
		html += h.getNaviLink(n...)
	}
	return fmt.Sprintf(`<nav>%s</nav>`, html)
}

func (h *NarrativePageHtml) getNaviLink(vals ...string) string {
	return fmt.Sprintf(`<a rel="%s" title="%s" href="%s">%s</a>`, vals[0], vals[1], vals[2], vals[3])
}

func (h *NarrativePageHtml) getMetaHtml() string {
	ms := h.p.GetMeta()
	html := h.HTML.getMetaHtml()
	for _, m := range ms {
		html += h.getHeaderLink(m...)
	}

	return html
}

func (h *NarrativePageHtml) getHeaderHtml() string {
	hl := h.getHeadline(h.p.Title())
	s := config.Servedrootpath()
	return fmt.Sprintf(`
	<a href="%s" class="home"><!--DevAbo.de--></a>
    <a href="%s/2013/08/01/a-step-in-the-dark/" class="orange">New Reader? Start here!</a>
	%s`, s, s, hl)
}

func (h *NarrativePageHtml) getDisqus() string {
	title := h.p.Title()
	url := h.getDisqusUrl()
	identifier := h.getDisqusIdentifier()
	disq := fmt.Sprintf(disqus_universal_code, title, url, identifier)
	return disq
}

func (h *NarrativePageHtml) getDisqusIdentifier() string {
	if len(h.p.DisqusIdentifier()) > 0 {
		return h.p.DisqusIdentifier()
	}
	return h.p.Path()
}

func (h *NarrativePageHtml) getDisqusUrl() string {
	return h.p.Path() + "/"
}

func DateNow() string {
	date := time.Now()
	return date.Format(time.RFC1123Z)
}
