package fs

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
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
	o.writeArchive()
}

func (o *Output) writeNarrativePages() {
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
}

func (o *Output) getImageAsBase64(p *page.Page) {
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
}

func (o *Output) writeArchive() {
	list := []string{}
	for _, p := range o.comic.GetPages() {
		o.getImageAsBase64(p)
		list = append(list, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, p.Path(), p.Title()))

	}

	arc := fmt.Sprintf("<ul>%s</ul>", strings.Join(list, "\n"))
	ah := NewArchiveHtml(arc)
	log.Println(ah.getContent())
	o.writeStringToFS(config.Rootpath()+"/archive.html", ah.writePage())
}

func (o *Output) writeCss() {
	p := config.Rootpath() + "/css"
	o.prepareFileSystem(p)
	fp := p + "/style.css"
	o.writeStringToFS(fp, css)
}

func (o *Output) writePageToFileSystem(p *page.Page) {
	absPath := config.Rootpath() + p.FSPath()
	o.prepareFileSystem(absPath)

	h := NewNarrativePageHtml(p)
	o.writeStringToFS(absPath+"/index.html", h.writePage())
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
	return fmt.Sprintf(`<a href="http://twitter.com/devabo_de">Twitter</a>
	<a href="%s/about.html">About</a>
	<a href="%s/archive.html">Archive</a>
	<a href="%s/imprint.html">Imprint / Impressum</a>
	`, p, p, p)
}

func (html *HTML) getCssLink() string {
	path := config.Servedrootpath() + "/css/style.css?version=" + html.version()
	format := `<link rel="stylesheet" href="%s" type="text/css">`
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
	css := html.getCssLink()
	meta := html.getMetaHtml()
	navi := html.getNaviHtml()
	title := html.getTitle()
	footerNavi := html.getFooterNavi()
	content := html.getContent()
	header := html.getHeaderHtml()
	disqus := ""
	year := time.Now().Year()
	return fmt.Sprintf(htmlFormat, title, meta, css, header, content, navi, disqus, year, footerNavi)
}

type ArchiveHtml struct {
	HTML
	content string
}

func NewArchiveHtml(content string) *ArchiveHtml {
	return &ArchiveHtml{HTML{}, content}
}

func (ah *ArchiveHtml) getContent() string {
	return ah.content
}
func (ah *ArchiveHtml) writePage() string {
	css := ah.getCssLink()
	meta := ah.getMetaHtml()
	navi := ah.getNaviHtml()
	title := ah.getTitle()
	footerNavi := ah.getFooterNavi()
	content := ah.getContent()
	header := ah.getHeaderHtml()
	disqus := ""
	year := time.Now().Year()
	return fmt.Sprintf(htmlFormat, title, meta, css, header, content, navi, disqus, year, footerNavi)
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
	css := h.getCssLink()
	meta := h.getMetaHtml()
	navi := h.getNaviHtml()
	title := h.p.Title()
	footerNavi := h.getFooterNavi()
	content := h.getContent()
	header := h.getHeaderHtml()
	disqus := h.getDisqus()
	year := time.Now().Year()
	return fmt.Sprintf(htmlFormat, title, meta, css, header, content, navi, disqus, year, footerNavi)
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
	html := ""
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

const css = `
.copyright,
header {
	width: 800px;
	margin: 0 auto;
}

.copyright{
	margin-top: 30px;
}

h3 {
	font-family: Arial Black;
	text-align: left;
	text-transform: uppercase;
}

header .home {
    display: block;
    line-height: 80px;
    background: url(https://devabo.de/wp-content/themes/drewing2012/header_devabo_de.png) no-repeat 0px -0px;
    height: 30px;
    width: 800px;
    text-align: left;
    color: #000;
    margin-bottom: 0px;
	margin-top: 0;
    background-color: transparent;
}
header .orange {
	display: block;
    height: 2.2em;
    background-color: #FF8800;
    color: #FFFFFF;
    line-height: 1em;
    padding: 0.5em;
    box-sizing: border-box;
	width: 100%;
    font-size: 24px;
	font-family: Arial Black;
	text-transform: uppercase;
    text-decoration: underline;
	margin-bottom: 1rem;
}

body {
	text-align: center;
	margin: 0;
	padding: 0;
	border: 0;
	font-family: Arial, Helvetica, sans-serif;
}

#disqus_thread,
main {
	width: 800px;
	margin: 0 auto;
}

footer {
	position: fixed;
	bottom: 0;
	width: 100%;
	text-align: center;
	z-index: 100;
}

footer nav {
	border-top: 1px solid black;
	position: relative;
	background-color: white;
	min-height: 45px;
	width: 800px;
	margin: 0 auto;
}

nav a {
	font-family: Arial Black;
	color: black;
	text-decoration: none;
	height: 100%;
	display: inline-block;
	padding: 10px;
	text-transform: uppercase;
}

.spacer {
	height: 80px;
}

`
const imageWrapperFormat = `<a href="%s" rel="next" title="%s">%s</a>`
const navWrapperFormat = `<nav>%s</nav>`
const htmlFormat = `<!DOCTYPE html>
<html>
	<head>
		<title>DevAbo.de | Graphic Novel | %s</title>
		%s
		%s
	</head>
	<body>
		<header>
%s
		</header>
		<main>
			%s
			%s
		</main>
		%s
		<div class="copyright">
		All content including but not limited to the art, characters, story, website design & graphics are © copyright 2013-%d Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.
		</div>
		<div class="spacer"></div>
		<footer><nav>%s</nav></footer>
	</body>
</html>
`

const disqus_universal_code = `
<div id="disqus_thread"></div>
<script type="text/javascript">
var disqus_title = "%s";
var disqus_url = 'https://DevAbo.de%s';
var disqus_identifier = '%s';
var disqus_container_id = 'disqus_thread';
var disqus_shortname = 'devabode';
var disqus_config_custom = window.disqus_config;
var disqus_config = function () {
    this.language = '';
        this.callbacks.onReady.push(function () {
        // sync comments in the background so we don't block the page
        var script = document.createElement('script');
        script.async = true;
        script.src = '?cf_action=sync_comments&post_id=1235';
        var firstScript = document.getElementsByTagName('script')[0];
        firstScript.parentNode.insertBefore(script, firstScript);
    });
    if (disqus_config_custom) {
        disqus_config_custom.call(this);
    }
};
(function() {
    var dsq = document.createElement('script');
	dsq.type = 'text/javascript';
    dsq.async = true;
	dsq.src = 'https://' + disqus_shortname + '.disqus.com/embed.js';
    (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
})();
</script>
`
