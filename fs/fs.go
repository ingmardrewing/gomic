package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
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
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
	o.writeCss()
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

	h := NewHtml(p)
	o.writeStringToFS(absPath+"/index.html", h.writePage())
}

func (o *Output) writeStringToFS(absPath string, html string) {
	log.Println("writing html to filesystem: ", absPath)
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

type Html struct {
	p                                           *page.Page
	title, meta, csslink, img, navi, footerNavi string
}

func NewHtml(p *page.Page) *Html {
	return &Html{p, "", "", "", "", "", ""}
}

func (h *Html) getFooterNavi() string {
	return `<a href="http://twitter.com/devabo_de">Twitter</a>
	<a href="/about.html">About</a>
	<a href="/archive.html">Archive</a>
	<a href="/imprint.html">Imprint / Impressum</a>
	`
}

func (h *Html) version() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02dT%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (h *Html) getCssLink() string {
	path := config.Servedrootpath() + "/css/style.css?version=" + h.version()
	format := `<link rel="stylesheet" href="%s" type="text/css">`
	return fmt.Sprintf(format, path)
}

func (h *Html) writePage() string {
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

func (h *Html) getContent() string {
	f := `<img src="%s" width="800" height="1334" alt="">`
	html := fmt.Sprintf(f, h.p.Img())
	if !h.p.IsLast() {
		html = fmt.Sprintf(`<a href="%s">%s</a>`, h.p.UrlToNext(), html)
	}
	return html
}

func (h *Html) getNaviHtml() string {
	ns := h.p.GetNavi()
	html := ""
	for _, n := range ns {
		html += h.getNaviLink(n...)
	}
	return fmt.Sprintf(`<nav>%s</nav>`, html)
}

func (h *Html) getNaviLink(vals ...string) string {
	return fmt.Sprintf(`<a rel="%s" title="%s" href="%s">%s</a>`, vals[0], vals[1], vals[2], vals[3])
}

func (h *Html) getMetaHtml() string {
	ms := h.p.GetMeta()
	html := ""
	for _, m := range ms {
		html += h.getHeaderLink(m...)
	}
	return html
}

func (h *Html) getHeaderLink(vals ...string) string {
	l := fmt.Sprintf(`<link rel="%s" title="%s" href="%s">`, vals[0], vals[1], vals[2])
	return l
}

func (h *Html) getHeadline() string {
	return fmt.Sprintf(`<h3>%s</h3>`, h.p.Title())
}

func (h *Html) getHeaderHtml() string {
	hl := h.getHeadline()
	return fmt.Sprintf(`
	<a href="https://DevAbo.de/" class="home"><!--DevAbo.de--></a>
    <a href="https://devabo.de/2013/08/01/a-step-in-the-dark/" class="orange">New Reader? Start here!</a>
	%s`, hl)
}

func (h *Html) getDisqus() string {
	title := h.p.Title()
	url := h.getDisqusUrl()
	identifier := h.getDisqusIdentifier()
	disq := fmt.Sprintf(disqus_universal_code, title, url, identifier)
	return disq
}

func (h *Html) getDisqusIdentifier() string {
	if len(h.p.DisqusIdentifier()) > 0 {
		return h.p.DisqusIdentifier()
	}
	return h.p.Path()
}

func (h *Html) getDisqusUrl() string {
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
