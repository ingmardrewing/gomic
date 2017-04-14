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
	o.writeStringToFS(config.Rootpath()+"/about.html", ah.writePage())
}

func (o *Output) writeImprint() {
	ah := NewDataHtml(imprint, config.Servedrootpath()+"/imprint.html")
	o.writeStringToFS(config.Rootpath()+"/imprint.html", ah.writePage())
}

func (o *Output) writeRss() {
	rss := o.Rss()
	o.writeStringToFS(config.Rootpath()+"/feed/rss.xml", rss)
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
	o.writeStringToFS(config.Rootpath()+"/archive.html", ah.writePage())
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
	<a href="%s/rss.xml">RSS</a>
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
	css_js := html.getCssLink() + html.getJsLink()
	meta := html.getMetaHtml()
	navi := html.getNaviHtml()
	title := html.getTitle()
	footerNavi := html.getFooterNavi()
	content := html.getContent()
	header := html.getHeaderHtml()
	disqus := ""
	year := time.Now().Year()
	canonicalLink := ""
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css_js, header, content, navi, disqus, year, footerNavi)
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
func (ah *DataHtml) writePage() string {
	css_js := ah.getCssLink() + ah.getJsLink()
	meta := ah.getMetaHtml()
	navi := ah.getNaviHtml()
	title := ah.getTitle()
	footerNavi := ah.getFooterNavi()
	content := ah.getContent()
	header := ah.getHeaderHtml()
	disqus := ""
	canonicalLink := fmt.Sprintf(`<link rel="canonical" href="%s">`, ah.url)
	year := time.Now().Year()
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css_js, header, content, navi, disqus, year, footerNavi)
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
	css_js := h.getCssLink() + h.getJsLink()
	meta := h.getMetaHtml()
	navi := h.getNaviHtml()
	title := h.p.Title()
	footerNavi := h.getFooterNavi()
	content := h.getContent()
	header := h.getHeaderHtml()
	disqus := h.getDisqus()
	year := time.Now().Year()
	canonicalLink := fmt.Sprintf(`<link rel="canonical" href="%s">`, h.p.Path())
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css_js, header, content, navi, disqus, year, footerNavi)
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

const css = `
.copyright,
header {
	width: 800px;
	margin: 0 auto;
}

ul.archive {
	list-style-type: none;
}

ul.archive li {
	display: inline-block;
	margin: 10px;
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
    background: url(https://devabo.de/imgs/header_devabo_de.png) no-repeat 0px -0px;
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

#cookie-law-info-bar {
	font-size: 10pt;
	margin: 0 auto;
	padding: 5px 0;
	position: fixed;
	top: 0;
	left: 0;
	text-align: center;
	width: 100%;
	z-index: 9999;
	background-color: white;
	border: 1px solid black;
}
#cookie-law-info-again {
	font-size: 10pt;
	margin: 0;
	padding: 2px 10px;
	text-align: center;
	z-index: 9999;
	cursor: pointer;
}
#cookie-law-info-bar span {
	vertical-align: middle;
}
/** Buttons (http://papermashup.com/demos/css-buttons) */
.cli-plugin-button, .cli-plugin-button:visited {
	display: inline-block;
	padding: 5px 10px 6px;
	color: #fff;
	text-decoration: none;
	-moz-border-radius: 6px;
	-webkit-border-radius: 6px;
	-moz-box-shadow: 0 1px 3px rgba(0,0,0,0.6);
	-webkit-box-shadow: 0 1px 3px rgba(0,0,0,0.6);
	text-shadow: 0 -1px 1px rgba(0,0,0,0.25);
	border-bottom: 1px solid rgba(0,0,0,0.25);
	position: relative;
	cursor: pointer;
	margin: auto 10px;
}
.cli-plugin-button:hover {
	background-color: #111;
	color: #fff;
}
.cli-plugin-button:active {
	top: 1px;
}
.small.cli-plugin-button, .small.cli-plugin-button:visited {
	font-size: 11px;
}
.cli-plugin-button, .cli-plugin-button:visited,
	.medium.cli-plugin-button, .medium.cli-plugin-button:visited {
	font-size: 13px;
	font-weight: bold;
	line-height: 1;
	text-shadow: 0 -1px 1px rgba(0,0,0,0.25);
}
.large.cli-plugin-button, .large.cli-plugin-button:visited {
	font-size: 14px;
	padding: 8px 14px 9px;
}
.super.cli-plugin-button, .super.cli-plugin-button:visited {
	font-size: 34px;
	padding: 8px 14px 9px;
}
.pink.cli-plugin-button, .magenta.cli-plugin-button:visited {
	background-color: #e22092;
}
.pink.cli-plugin-button:hover {
	background-color: #c81e82;
}
.green.cli-plugin-button, .green.cli-plugin-button:visited {
	background-color: #91bd09;
}
.green.cli-plugin-button:hover {
	background-color: #749a02;
}
.red.cli-plugin-button, .red.cli-plugin-button:visited {
	background-color: #e62727;
}
.red.cli-plugin-button:hover {
	background-color: #cf2525;
}
.orange.cli-plugin-button, .orange.cli-plugin-button:visited {
	background-color: #ff5c00;
}
.orange.cli-plugin-button:hover {
	background-color: #d45500;
}
.blue.cli-plugin-button, .blue.cli-plugin-button:visited {
	background-color: #2981e4;
}
.blue.cli-plugin-button:hover {
	background-color: #2575cf;
}
.yellow.cli-plugin-button, .yellow.cli-plugin-button:visited {
	background-color: #ffb515;
}
.yellow.cli-plugin-button:hover {
	background-color: #fc9200;
}

`
const imageWrapperFormat = `<a href="%s" rel="next" title="%s">%s</a>`
const navWrapperFormat = `<nav>%s</nav>`

const htmlFormat = `<!DOCTYPE html>
<html itemscope itemtype="http://schema.org/Article" lang="en" manifest="/cache.manifest" >
	<head>
		<meta http-equiv="content-type" content="text/html;charset=UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="robots" content="index,follow">
		<meta name="author" content="Ingmar Drewing"> 
		<meta name="publisher" content="Ingmar Drewing"> 
		<meta name="keywords" content="web comic, comic, cartoon, sci fi, satire, parody, science fiction, action, software industry, pulp, nerd, geek"> 
		<meta name="DC.Subject" content="web comic, comic, cartoon, sci fi, science fiction, satire, parody action, software industry"> 
		<meta name="page-topic" content="Science Fiction Web-Comic">
		<link rel="dns-prefetch" href="https://DevAbo.de">

		<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
		<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
		<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
		<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
		<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
		<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
		<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
		<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
		<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
		<link rel="icon" type="image/png" sizes="192x192"  href="/android-icon-192x192.png">
		<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
		<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		<link rel="manifest" href="/manifest.json">
		<meta name="msapplication-TileColor" content="#ffffff">
		<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
		<meta name="theme-color" content="#ffffff">
		<meta name="description" content=" " /><!-- 155 chars -->

		<!-- 
		Schema.org markup for Google+ 
		<meta itemprop="name" content="The Name or Title Here">
		<meta itemprop="description" content="This is the page description">
		<meta itemprop="image" content="http://www.example.com/image.jpg">

 Twitter Card data 
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:site" content="@publisher_handle">
<meta name="twitter:title" content="Page Title">
<meta name="twitter:description" content="Page description less than 200 characters">
<meta name="twitter:creator" content="@author_handle">
-->

<!-- Twitter summary card with large image must be at least 280x150px 
<meta name="twitter:image:src" content="http://www.example.com/image.jpg">
-->

<!-- 
Open Graph data
<meta property="og:title" content="Title Here" />
<meta property="og:type" content="article" />
<meta property="og:url" content="http://www.example.com/" />
<meta property="og:image" content="http://example.com/image.jpg" />
<meta property="og:description" content="Description Here" />
<meta property="og:site_name" content="Site Name, i.e. Moz" />
<meta property="article:published_time" content="2013-09-17T05:59:00+01:00" />
<meta property="article:modified_time" content="2013-09-16T19:08:47+01:00" />
<meta property="article:section" content="Article Section" />
<meta property="article:tag" content="Article Tag" />
<meta property="fb:admins" content="Facebook numberic ID" />

-->

		%s
		<title>test DevAbo.de | Graphic Novel | %s</title>
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
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
		All content including but not limited to the art, characters, story, website design & graphics are &copy; copyright 2013-%d Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.
		    <div id="cookie-law-info-bar"><span>This website uses cookies to improve your experience. We'll assume you're ok with this, but you can opt-out if you wish.<a href="#" id="cookie_action_close_header"  class="medium cli-plugin-button cli-plugin-main-button" >Accept</a> <a href="http://www.drewing.de/blog/impressum-imprint/" id="CONSTANT_OPEN_URL" target="_blank"  class="cli-plugin-main-link"  >Read More</a></span></div>		
		</div>
		<div class="spacer"></div>
		<footer><nav>%s</nav></footer>
	</body>
</html>
`

const disqus_universal_code = `
<div id="disqus_thread"></div>
<script type="text/javascript">
//<![CDATA[


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

//]]>
</script>
`

var analytics = `

<script type="text/javascript">
//<![CDATA[
var gaProperty = 'UUA-49679648-1';
var disableStr = 'ga-disable-' + gaProperty;
if (document.cookie.indexOf(disableStr + '=true') > -1) {
  window[disableStr] = true;
}
function gaOptout() {
  document.cookie = disableStr + '=true; expires=Thu, 31 Dec 2099 23:59:59 UTC; path=/';
  window[disableStr] = true;
}

  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-49679648-1', 'devabo.de');
  ga('set', 'anonymizeIp', true);
  ga('require', 'displayfeatures');
  ga('require', 'linkid', 'linkid.js');
  ga('send', 'pageview');

//]]>
</script>`

func (o *Output) RssItem(p *page.Page) string {
	title := p.Title()
	url := p.Path()
	pubDate := p.Date()
	act := p.Act()
	description := p.Title()
	content := fmt.Sprintf(`<img src="%s">`, p.ImgUrl())
	thumbnailUrl := p.ThumnailUrl()
	imageUrl := p.ImgUrl()
	imageName := p.ImageFilename()
	return fmt.Sprintf(rssItem, title, url, pubDate, act, url, description, content, thumbnailUrl, imageUrl, imageName, thumbnailUrl)
}

func (o *Output) RssItems() string {
	h := ""
	pgs := o.comic.GetPages()
	// last 10 pages
	l10 := pgs[len(pgs)-11:]

	// reverse splice
	for i := len(l10)/2 - 1; i >= 0; i-- {
		opp := len(l10) - 1 - i
		l10[i], l10[opp] = l10[opp], l10[i]
	}
	// generate rss for last 10 pages, reversed
	for _, p := range l10 {
		h += o.RssItem(p)
	}
	return h
}

func (o *Output) DateNow() string {
	date := time.Now()
	return date.Format(time.RFC1123Z)
}

func (o *Output) Rss() string {
	date := o.DateNow()
	items := o.RssItems()
	relSelf := config.Servedrootpath() + "/feed/rss.xml"
	return fmt.Sprintf(rss, relSelf, date, items)
}

var rss = `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wfw="http://wellformedweb.org/CommentAPI/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:atom="http://www.w3.org/2005/Atom"
	xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
	xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
	xmlns:media="http://search.yahoo.com/mrss/"
	>

<channel>
	<title>DevAbo.de</title>
    <image>
      <url>https://devabo.de/favicon-32x32.png</url>
      <title>DevAbo.de</title>
      <link>https://devabo.de</link>
      <width>32</width>
      <height>32</height>
      <description>A science-fiction webcomic about the lives of software developers in the far, funny and dystopian future</description>
    </image>
	<atom:link href="%s" rel="self" type="application/rss+xml" />
	<link>https://DevAbo.de</link>
	<description>A science-fiction webcomic about the lives of software developers in the far, funny and dystopian future</description>
	<lastBuildDate>%s</lastBuildDate>
	<language>en-US</language>
	<sy:updatePeriod>weekly</sy:updatePeriod>
	<sy:updateFrequency>1</sy:updateFrequency>
	<generator>https://github.com/ingmardrewing/gomic</generator>
%s
	</channel>
</rss>
`

var rssItem = `  <item>
    <title>%s</title>
    <link>%s</link>
    <pubDate>%s</pubDate>
    <dc:creator><![CDATA[Ingmar Drewing]]></dc:creator>
    <category><![CDATA[%s]]></category>
    <guid>%s/index.html</guid>
    <description><![CDATA[%s]]></description>
    <content:encoded><![CDATA[%s]]></content:encoded>

    <media:thumbnail url="%s" />
    <media:content url="%s" medium="image">
      <media:title type="html">%s</media:title>
      <media:thumbnail url="%s" />
    </media:content>
  </item>
`

var imprint = `
Angaben nach TDG:

Dieses Impressum gilt für die Website devabo.de, so wie z.B. die zugehörige Facebook-Page <a href="https://www.facebook.com/devabo.de">https://www.facebook.com/devabo.de</a> und das Facebook-Profil <a href="https://www.facebook.com/ingmar.drewing">https://www.facebook.com/ingmar.drewing</a> sowie die Google-Plus Seite <a href="https://plus.google.com/107755781256885973202/posts">https://plus.google.com/107755781256885973202/posts</a> der Twitter-Account unter <a href="https://twitter.com/ingmardrewing">https://twitter.com/ingmardrewing</a> und alle weiteren Profile und Websites von Ingmar Drewing, wie auch www.devabo.de .

Redaktionell verantwortlich ist:

Ingmar Drewing
(Dipl. Kommunkationsdesigner /FH /BRD)
Schulberg 8
65183 Wiesbaden

Telefon: 0173-3076520
E-Mail: ingmar-at-drewing-punkt-de

Newsletter
Diese Website bietet eine Newsletter-Funktion an. Heißt: Du kannst deine Mailadresse in das Newsletter-Feld oben eintragen und wirst dann von mir per E-Mail informiert, wenn eine neue Comic-Seite online geht. Ich versende den Newsletter ausschließlich an interessierte Leser, die den Versandt ausdrücklich bestellt haben. Für den Versandt verwende ich den Listenprovider MailChimp. MailChimp gehört The Rocket Science Group LLC, 512 Means Street, Ste 404 Atlanta, GA 30318. 

Wenn Du dich für meinen Newsletter registriert werden die Daten, die du bei der Newsletterregistrierung eingegeben hast (also die Mailadresse und ggf. noch Vor- und Zuname, sofern du die Felder freiwillig über deren Website ergänzt - da ist ein Link dazu in den Newsletter-Mails enthalten), dorthin übertragen und dort gespeichert. Nach der Anmeldung wird erst mal eine Bestätigungsmail an die von dir angegebene E-Mail-Adresse versandt, um die Bestellung zu verifizieren ("double opt-in").

Mailchimp hat ein breites Spektrum an Analyseinstrumenten, wie die Newsletter genutzt und geöffnet werden. Diese Analysen sind gruppenbezogen und werden von mir nicht verwendet um das Verhalten einzelner Individuen zu betrachten. MailChimp nutzt auch Google Analystics und bindet es eventuell in den verschickten Newsletter ein. Die verschickten Newsletter-E-Mails enthalten alle einen Link zum abbestellen, falls du den Newsletter nicht mehr erhalten möchtest.
    
Haftungsausschluss

1. Inhalt des Onlineangebotes
Der Autor übernimmt keinerlei Gewähr für die Aktualität, Korrektheit, Vollständigkeit oder Qualität der bereitgestellten Informationen. Haftungsansprüche gegen den Autor, welche sich auf Schäden materieller oder ideeller Art beziehen, die durch die Nutzung oder Nichtnutzung der dargebotenen Informationen bzw. durch die Nutzung fehlerhafter und unvollständiger Informationen verursacht wurden, sind grundsätzlich ausgeschlossen, sofern seitens des Autors kein nachweislich vorsätzliches oder grob fahrlässiges Verschulden vorliegt.
Alle Angebote sind freibleibend und unverbindlich. Der Autor behält es sich ausdrücklich vor, Teile der Seiten oder das gesamte Angebot ohne gesonderte Ankündigung zu verändern, zu ergänzen, zu löschen oder die Veröffentlichung zeitweise oder endgültig einzustellen.

2. Verweise und Links
Bei direkten oder indirekten Verweisen auf fremde Webseiten ("Hyperlinks"), die außerhalb des Verantwortungsbereiches des Autors liegen, würde eine Haftungsverpflichtung ausschließlich in dem Fall in Kraft treten, in dem der Autor von den Inhalten Kenntnis hat und es ihm technisch möglich und zumutbar wäre, die Nutzung im Falle rechtswidriger Inhalte zu verhindern.
Der Autor erklärt hiermit ausdrücklich, dass zum Zeitpunkt der Linksetzung keine illegalen Inhalte auf den zu verlinkenden Seiten erkennbar waren. Auf die aktuelle und zukünftige Gestaltung, die Inhalte oder die Urheberschaft der verlinkten/verknüpften Seiten hat der Autor keinerlei Einfluss. Deshalb distanziert er sich hiermit ausdrücklich von allen Inhalten aller verlinkten /verknüpften Seiten, die nach der Linksetzung verändert wurden. Diese Feststellung gilt für alle innerhalb des eigenen Internetangebotes gesetzten Links und Verweise sowie für Fremdeinträge in vom Autor eingerichteten Gästebüchern, Diskussionsforen, Linkverzeichnissen, Mailinglisten und in allen anderen Formen von Datenbanken, auf deren Inhalt externe Schreibzugriffe möglich sind. Für illegale, fehlerhafte oder unvollständige Inhalte und insbesondere für Schäden, die aus der Nutzung oder Nichtnutzung solcherart dargebotener Informationen entstehen, haftet allein der Anbieter der Seite, auf welche verwiesen wurde, nicht derjenige, der über Links auf die jeweilige Veröffentlichung lediglich verweist.

3. Urheber- und Kennzeichenrecht
Der Autor ist bestrebt, in allen Publikationen die Urheberrechte der verwendeten Bilder, Grafiken, Tondokumente, Videosequenzen und Texte zu beachten, von ihm selbst erstellte Bilder, Grafiken, Tondokumente, Videosequenzen und Texte zu nutzen oder auf lizenzfreie Grafiken, Tondokumente, Videosequenzen und Texte zurückzugreifen.
Alle innerhalb des Internetangebotes genannten und ggf. durch Dritte geschützten Marken- und Warenzeichen unterliegen uneingeschränkt den Bestimmungen des jeweils gültigen Kennzeichenrechts und den Besitzrechten der jeweiligen eingetragenen Eigentümer. Allein aufgrund der bloßen Nennung ist nicht der Schluss zu ziehen, dass Markenzeichen nicht durch Rechte Dritter geschützt sind!
Das Copyright für veröffentlichte, vom Autor selbst erstellte Objekte bleibt allein beim Autor der Seiten. Eine Vervielfältigung oder Verwendung solcher Grafiken, Tondokumente, Videosequenzen und Texte in anderen elektronischen oder gedruckten Publikationen ist ohne ausdrückliche Zustimmung des Autors nicht gestattet.

4. Datenschutz
Sofern innerhalb des Internetangebotes die Möglichkeit zur Eingabe persönlicher oder geschäftlicher Daten (Emailadressen, Namen, Anschriften) besteht, so erfolgt die Preisgabe dieser Daten seitens des Nutzers auf ausdrücklich freiwilliger Basis. Die Inanspruchnahme und Bezahlung aller angebotenen Dienste ist - soweit technisch möglich und zumutbar - auch ohne Angabe solcher Daten bzw. unter Angabe anonymisierter Daten oder eines Pseudonyms gestattet. Die Nutzung der im Rahmen des Impressums oder vergleichbarer Angaben veröffentlichten Kontaktdaten wie Postanschriften, Telefon- und Faxnummern sowie Emailadressen durch Dritte zur Übersendung von nicht ausdrücklich angeforderten Informationen ist nicht gestattet. Rechtliche Schritte gegen die Versender von sogenannten Spam-Mails bei Verstössen gegen dieses Verbot sind ausdrücklich vorbehalten.

5. Rechtswirksamkeit dieses Haftungsausschlusses
Dieser Haftungsausschluss ist als Teil des Internetangebotes zu betrachten, von dem aus auf diese Seite verwiesen wurde. Sofern Teile oder einzelne Formulierungen dieses Textes der geltenden Rechtslage nicht, nicht mehr oder nicht vollständig entsprechen sollten, bleiben die übrigen Teile des Dokumentes in ihrem Inhalt und ihrer Gültigkeit davon unberührt.

6. Google Analytics (Text übernommen von <a href="http://www.datenschutzbeauftragter-info.de">www.datenschutzbeauftragter-info.de</a>)
Diese Website benutzt Google Analytics, einen Webanalysedienst der Google Inc. („Google“). Google Analytics verwendet sog. „Cookies“, Textdateien, die auf Ihrem Computer gespeichert werden und die eine Analyse der Benutzung der Website durch Sie ermöglichen. Die durch den Cookie erzeugten Informationen über Ihre Benutzung dieser Website werden in der Regel an einen Server von Google in den USA übertragen und dort gespeichert. Im Falle der Aktivierung der IP-Anonymisierung auf dieser Website, wird Ihre IP-Adresse von Google jedoch innerhalb von Mitgliedstaaten der Europäischen Union oder in anderen Vertragsstaaten des Abkommens über den Europäischen Wirtschaftsraum zuvor gekürzt. Nur in Ausnahmefällen wird die volle IP-Adresse an einen Server von Google in den USA übertragen und dort gekürzt. Im Auftrag des Betreibers dieser Website wird Google diese Informationen benutzen, um Ihre Nutzung der Website auszuwerten, um Reports über die Websiteaktivitäten zusammenzustellen und um weitere mit der Websitenutzung und der Internetnutzung verbundene Dienstleistungen gegenüber dem Websitebetreiber zu erbringen. Die im Rahmen von Google Analytics von Ihrem Browser übermittelte IP-Adresse wird nicht mit anderen Daten von Google zusammengeführt. Sie können die Speicherung der Cookies durch eine entsprechende Einstellung Ihrer Browser-Software verhindern; wir weisen Sie jedoch darauf hin, dass Sie in diesem Fall gegebenenfalls nicht sämtliche Funktionen dieser Website vollumfänglich werden nutzen können. Sie können darüber hinaus die Erfassung der durch das Cookie erzeugten und auf Ihre Nutzung der Website bezogenen Daten (inkl. Ihrer IP-Adresse) an Google sowie die Verarbeitung dieser Daten durch Google verhindern, indem sie das unter dem folgenden Link (<a href="http://tools.google.com/dlpage/gaoptout?hl=de">http://tools.google.com/dlpage/gaoptout?hl=de</a>) verfügbare Browser-Plugin herunterladen und installieren.

Sie können die Erfassung durch Google Analytics verhindern, indem Sie auf folgenden Link klicken. Es wird ein Opt-Out-Cookie gesetzt, der die zukünftige Erfassung Ihrer Daten beim Besuch dieser Website verhindert:
<a href="javascript:gaOptout()">Google Analytics deaktivieren</a>

Nähere Informationen zu Nutzungsbedingungen und Datenschutz finden Sie unter <a href="http://www.google.com/analytics/terms/de.html">http://www.google.com/analytics/terms/de.html</a> bzw. unter <a href="http://www.google.com/intl/de/analytics/privacyoverview.html">http://www.google.com/intl/de/analytics/privacyoverview.html</a>. Wir weisen Sie darauf hin, dass auf dieser Website Google Analytics um den Code „gat._anonymizeIp();“ erweitert wurde, um eine anonymisierte Erfassung von IP-Adressen (sog. IP-Masking) zu gewährleisten.


Disclaimer

1. Content
The author reserves the right not to be responsible for the topicality, correctness, completeness or quality of the information provided. Liability claims regarding damage caused by the use of any information provided, including any kind of information which is incomplete or incorrect,will therefore be rejected.
All offers are not-binding and without obligation. Parts of the pages or the complete publication including all offers and information might be extended, changed or partly or completely deleted by the author without separate announcement.

2. Referrals and links
The author is not responsible for any contents linked or referred to from his pages - unless he has full knowledge of illegal contents and would be able to prevent the visitors of his site fromviewing those pages. If any damage occurs by the use of information presented there, only the author of the respective pages might be liable, not the one who has linked to these pages. Furthermore the author is not liable for any postings or messages published by users of discussion boards, guestbooks or mailinglists provided on his page.

3. Copyright
The author intended not to use any copyrighted material for the publication or, if not possible, to indicate the copyright of the respective object.
The copyright for any material created by the author is reserved. Any duplication or use of objects such as images, diagrams, sounds or texts in other electronic or printed publications is not permitted without the author's agreement.

4. Privacy policy
If the opportunity for the input of personal or business data (email addresses, name, addresses) is given, the input of these data takes place voluntarily. The use and payment of all offered services are permitted - if and so far technically possible and reasonable - without specification of any personal data or under specification of anonymized data or an alias. The use of published postal addresses, telephone or fax numbers and email addresses for marketing purposes is prohibited, offenders sending unwanted spam messages will be punished.

5. Legal validity of this disclaimer
This disclaimer is to be regarded as part of the internet publication which you were referred from. If sections or individual terms of this statement are not legal or correct, the content or validity of the other parts remain uninfluenced by this fact.

6. Google Analytics (Text by <a href="http://www.datenschutzbeauftragter-info.de">www.datenschutzbeauftragter-info.de</a>)
This website uses Google Analytics, a web analytics service provided by Google, Inc. (“Google”).  Google Analytics uses “cookies”, which are text files placed on your computer, to help the website analyze how users use the site. The information generated by the cookie about your use of the website (including your IP address) will be transmitted to and stored by Google on servers in the United States.  In case of activation of the IP anonymization, Google will truncate/anonymize the last octet of the IP address for Member States of the European Union as well as for other parties to the Agreement on the European Economic Area.  Only in exceptional cases, the full IP address is sent to and shortened by Google servers in the USA.  On behalf of the website provider Google will use this information for the purpose of evaluating your use of the website, compiling reports on website activity for website operators and providing other services relating to website activity and internet usage to the website provider.  Google will not associate your IP address with any other data held by Google.  You may refuse the use of cookies by selecting the appropriate settings on your browser. However, please note that if you do this, you may not be able to use the full functionality of this website.  Furthermore you can prevent Google’s collection and use of data (cookies and IP address) by downloading and installing the browser plug-in available under <a href="https://tools.google.com/dlpage/gaoptout?hl=en-GB">https://tools.google.com/dlpage/gaoptout?hl=en-GB</a>.
You can refuse the use of Google Analytics by clicking on the following link. An opt-out cookie will be set on the computer, which prevents the future collection of your data when visiting this website:
<a href="javascript:gaOptout()">Disable Google Analytics</a>
Further information concerning the terms and conditions of use and data privacy can be found at <a href="http://www.google.com/analytics/terms/gb.html">http://www.google.com/analytics/terms/gb.html</a> or at <a href="http://www.google.com/intl/en_uk/analytics/privacyoverview.html">http://www.google.com/intl/en_uk/analytics/privacyoverview.html</a>.  Please note that on this website, Google Analytics code is supplemented by “gat._anonymizeIp();” to ensure an anonymized collection of IP addresses (so called IP-masking).
`

var about = `
devabo.de is a software developers fever dream, a nightmarish science-fiction comic by Ingmar Drewing. You can find the associated facebook-page at <a href="https://www.facebook.com/www.devabo.de">https://www.facebook.com/www.devabo.de</a> and the associated twitter account at <a href="twitter.com/#!/devabo_de">twitter.com/#!/devabo_de</a>.<br />
DevAbo.de will be updated every 1st and 15th day of every month, though I am trying right now to speed the production up to a weekly release cycle (but that's still beta).


<h3><a name="Bram">Bram</a></h3>
<img src="http://DevAbo.de/wp-content/uploads/2014/10/DevAbode_Bram.png" alt="DevAbo.de character Bram, comic, graphic novel, web comic, sci-fi, science-fiction" width="300" height="340" class="alignnone size-full wp-image-481" style="float:left; margin-right:25px; margin-bottom:15px;" />Bram has spent over a millennium in cryo stasis. <a href="#Ada">Ada</a> found him in the ancient ruins and ended his cryostatic slumber out of curiosity (and the possibility that he might have taken something of value into the cryo capsule). He woke up healthy, though a big part of his episodic memory is lost to him and resurfaces partially and slowly. <br /><br />
Some parts of his past that came back to him showed that he was some kind of <a href="http://devabo.de/2014/04/01/flashback/">technical officer</a>. He first didn't recall the aggressor he was fighting against. The memory of this came back to him while he was teaching Ada how to get in touch in with the calculating space and she <a href="http://devabo.de/2014/07/15/22-backup/">accidentally changed parts of his memory</a>.<br /><br />
Bram was about 35 years old when he was put into cryo slumber. He is still failing to remember the reason and circumstances of him being put into cryostasis, though it's likely that it has something to do with the war he was fighting in the past &hellip;<br /><br /><br /><br />


<h3><a name="Ada">Ada</a></h3>
<img src="http://DevAbo.de/wp-content/uploads/2014/10/DevAbode_Ada.png" alt="DevAbode.de Ada, scifi, science-fiction, character, comic, webcomic, graphic novel" width="300" height="341" class="alignnone size-full wp-image-461" style="float:right; margin-left:25px; margin-bottom:15px;" />Ada is a developer of the abode as well as an elite fighter. At the beginning of the story she is 27 years old.<br />
On routine checks along the outer defense perimeters of the abode she found entries to some rather well preserved buildings of the ancients and started to sell the artifacts she found there to <a href="#MasterBranch">Master Branch</a>. The business relationship developed and she regularly helped to retrieve artifacts for Master Branch. The business already brought her into conflict with her superiors and though she usually tries to keep out of trouble with the administration of the abode, she sneaked out into the ruins occasionally to "check for some strange client activity", as she put it in her report.<br />
Apart from this she takes her duty very seriously and is a good comrade. If a friend of her is in danger she's more than ready to risk her own life to free him. And she expects the same behaviour from everyone of her comrades.
<br /><br /><br /><br />

<h3><a name="MasterBranch">Master Branch</a></h3>
<img src="http://DevAbo.de/wp-content/uploads/2013/07/DevAbode_master_branch.png" alt="DevAbo.de character: Master Branch" width="300" height="340" class="alignnone size-full wp-image-450" style="float:left; margin-right:25px; margin-bottom:15px;" /> Master Branch is a thrirty years old JMonk and, like all of these pious people, believes strongly in <a href="http://en.wikipedia.org/wiki/Type_system#Static_type-checking">static typing</a>. Since that faith is mercilessly tested every time reality interferes with their believe system, the JMonks are also on the lookout for a sign from a higher power. They have a prophecy that one day a man would come, a Messiah, who will bring them true productivity. But, until this comes true, their only joy will be the incredible beauty of generics and jverbosity&trade;.<br /><br />
However, they still managed to create a machine that emenates fields of unproductivity. Though the specs didn't say the machine would do this, many of the JMonks hope that the machine might perhaps prove useful after all, one day.
Because of the difficulties mentioned above, Master Branch made a deal with a developer, <a href="#Ada">Ada</a>, whom he bid to go and search the ancient ruins for useful artifacts. He hoped to reverse engineer the artifacts and maybe find a way to become productive. Ada found and delivered several artifacts, which seemed interesting. Unfortunately they didn't reveal their usefulness yet. <br />
The most peculiar thing she found in the ruins she didn't deliver to the monks at all ...<br /><br /><br /><br />

<h3>Clients</h3>
Clients are ruthless and aggressive and also rather dim. That's making them less of a threat, as long as you are sufficiently quick witted. Nevertheless a greater pack of them can be quite distressing for a developer or consultant.
<br/><br/><br/><br/>


<h2>The Author</h2>

If you are interested in the other things I am drawing and writing you'll find some fragments on my <a href="http://www.drewing.de/blog">blog</a>. A word of warning: this blog contains material some people might consider nsfw.

`

var js = `
jQuery(document).ready(function() {
function cli_show_cookiebar(p) {
	/* plugin version 1.5.3 */
	var Cookie = {
		set: function(name,value,days) {
			if (days) {
				var date = new Date();
				date.setTime(date.getTime()+(days*24*60*60*1000));
				var expires = "; expires="+date.toGMTString();
			}
			else var expires = "";
			document.cookie = name+"="+value+expires+"; path=/";
		},
		read: function(name) {
			var nameEQ = name + "=";
			var ca = document.cookie.split(';');
			for(var i=0;i < ca.length;i++) {
				var c = ca[i];
				while (c.charAt(0)==' ') {
					c = c.substring(1,c.length);
				}
				if (c.indexOf(nameEQ) === 0) {
					return c.substring(nameEQ.length,c.length);
				}
			}
			return null;
		},
		erase: function(name) {
			this.set(name,"",-1);
		},
		exists: function(name) {
			return (this.read(name) !== null);
		}
	};
	
	var ACCEPT_COOKIE_NAME = 'viewed_cookie_policy',
		ACCEPT_COOKIE_EXPIRE = 365,
		json_payload = p.settings;
	
	if (typeof JSON.parse !== "function") {
		console.log("CookieLawInfo requires JSON.parse but your browser doesn't support it");
		return;
	}
	var settings = JSON.parse(json_payload);

	var cached_header = jQuery(settings.notify_div_id),
		cached_showagain_tab = jQuery(settings.showagain_div_id),
		btn_accept = jQuery('#cookie_hdr_accept'),
		btn_decline = jQuery('#cookie_hdr_decline'),
		btn_moreinfo = jQuery('#cookie_hdr_moreinfo'),
		btn_settings = jQuery('#cookie_hdr_settings');
	
	cached_header.hide();
	if ( !settings.showagain_tab ) {
		cached_showagain_tab.hide();
	}
	
	var hdr_args = { };
	
	var showagain_args = { };
	cached_header.css( hdr_args );
	cached_showagain_tab.css( showagain_args );
	
	if (!Cookie.exists(ACCEPT_COOKIE_NAME)) {
		displayHeader();
	}
	else {
		cached_header.hide();
	}
	
	if ( settings.show_once_yn ) {
		setTimeout(close_header, settings.show_once);
	}
	function close_header() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);
		hideHeader();
	}
	
	var main_button = jQuery('.cli-plugin-main-button');
	main_button.css( 'color', settings.button_1_link_colour );
	
	if ( settings.button_1_as_button ) {
		main_button.css('background-color', settings.button_1_button_colour);
		
		main_button.hover(function() {
			jQuery(this).css('background-color', settings.button_1_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_1_button_colour);
		});
	}
	var main_link = jQuery('.cli-plugin-main-link');
	main_link.css( 'color', settings.button_2_link_colour );
	
	if ( settings.button_2_as_button ) {
		main_link.css('background-color', settings.button_2_button_colour);
		
		main_link.hover(function() {
			jQuery(this).css('background-color', settings.button_2_button_hover);
		},
		function() {
			jQuery(this).css('background-color', settings.button_2_button_colour);
		});
	}
	
	cached_showagain_tab.click(function(e) {	
		e.preventDefault();
		cached_showagain_tab.slideUp(settings.animate_speed_hide, function slideShow() {
			cached_header.slideDown(settings.animate_speed_show);
		});
	});
	
	jQuery("#cookielawinfo-cookie-delete").click(function() {
		Cookie.erase(ACCEPT_COOKIE_NAME);
		return false;
	});
console.log("XXX");
console.log(jQuery("#cookie_action_close_header"));
console.log("XXX");
	jQuery("#cookie_action_close_header").click(function(e) {
		console.log('cookie_action_close_header clicked');
		e.preventDefault();
		accept_close();
	});

	function accept_close() {
		Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);
		
		if (settings.notify_animate_hide) {
			cached_header.slideUp(settings.animate_speed_hide);
		}
		else {
			cached_header.hide();
		}
		cached_showagain_tab.slideDown(settings.animate_speed_show);
		return false;
	}

	function closeOnScroll() {
		if (window.pageYOffset > 100 && !Cookie.read(ACCEPT_COOKIE_NAME)) {
			accept_close();
			if (settings.scroll_close_reload === true) {
				location.reload();
			}
			window.removeEventListener("scroll", closeOnScroll, false);
		}
	}
	if (settings.scroll_close === true) {
		window.addEventListener("scroll", closeOnScroll, false);
	}
	
	function displayHeader() {
		if (settings.notify_animate_show) {
			cached_header.slideDown(settings.animate_speed_show);
		}
		else {
			cached_header.show();
		}
		cached_showagain_tab.hide();
	}
	function hideHeader() {
		if (settings.notify_animate_show) {
			cached_showagain_tab.slideDown(settings.animate_speed_show);
		}
		else {
			cached_showagain_tab.show();
		}
		cached_header.slideUp(settings.animate_speed_show);
	}
};
function l1hs(str){if(str.charAt(0)=="#"){str=str.substring(1,str.length);}else{return "#"+str;}return l1hs(str);}

				cli_show_cookiebar({
					settings: '{"animate_speed_hide":"500","animate_speed_show":"500","background":"#fff","border":"#444","border_on":true,"button_1_button_colour":"#000","button_1_button_hover":"#000000","button_1_link_colour":"#fff","button_1_as_button":true,"button_2_button_colour":"#333","button_2_button_hover":"#292929","button_2_link_colour":"#444","button_2_as_button":false,"font_family":"inherit","header_fix":false,"notify_animate_hide":true,"notify_animate_show":false,"notify_div_id":"#cookie-law-info-bar","notify_position_horizontal":"right","notify_position_vertical":"bottom","scroll_close":false,"scroll_close_reload":false,"showagain_tab":false,"showagain_background":"#fff","showagain_border":"#000","showagain_div_id":"#cookie-law-info-again","showagain_x_position":"100px","text":"#000","show_once_yn":false,"show_once":"10000"}'
				});
			});
`
