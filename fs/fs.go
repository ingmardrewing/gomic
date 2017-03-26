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
	o.writeArchive()
	o.writeRss()
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
	ah := NewArchiveHtml(arc, config.Servedrootpath()+"/archive.html")
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
	`, p, p, p)
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
	canonicalLink := ""
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css, header, content, navi, disqus, year, footerNavi)
}

type ArchiveHtml struct {
	HTML
	content string
	url     string
}

func NewArchiveHtml(content string, url string) *ArchiveHtml {
	return &ArchiveHtml{HTML{}, content, url}
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
	canonicalLink := fmt.Sprintf(`<link rel="canonical" href="%s">`, ah.url)
	year := time.Now().Year()
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css, header, content, navi, disqus, year, footerNavi)
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
	canonicalLink := fmt.Sprintf(`<link rel="canonical" href="%s">`, h.p.Path())
	return fmt.Sprintf(htmlFormat, canonicalLink, title, meta, css, header, content, navi, disqus, year, footerNavi)
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
	<head profile="http://gmpg.org/xfn/11">
		<meta http-equiv="imagetoolbar" content="no">
		<meta http-equiv="content-type" content="text/html;charset=UTF-8">
		<meta http-equiv="Language" content="en">
		<meta http-equiv="Content-Language" content="en">
		<meta http-equiv="cache-control" content="Private">
		<meta http-equiv="pragma" content="no-cache">
		<meta http-equiv="expires" content="0">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="robots" content="index,follow">
		<meta name="author" content="Ingmar Drewing"> 
		<meta name="publisher" content="Ingmar Drewing"> 
		<meta name="keywords" content="web comic, comic, cartoon, sci fi, satire, parody, science fiction, action, software industry, pulp, nerd, geek"> 
		<meta name="DC.Subject" content="web comic, comic, cartoon, sci fi, science fiction, satire, parody action, software industry"> 
		<meta name="page-topic" content="Science Fiction Web-Comic">
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
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
		%s


		<!-- TODO:

		<link rel="shortcut icon" href="https://DevAbo.de/wp-content/themes/drewing2012/favicon.ico">
		<link rel="alternate" type="application/rss+xml" title="DevAbo.de » Feed" href="https://DevAbo.de/feed/">
		<script type="text/javascript" src="https://DevAbo.de/wp-content/plugins/cookie-law-info/js/cookielawinfo.js?ver=1.5.3"></script>
		-->

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
		All content including but not limited to the art, characters, story, website design & graphics are &copy; copyright 2013-%d Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.
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

var analytics = `
<script>
var gaProperty = 'UUA-49679648-1';
var disableStr = 'ga-disable-' + gaProperty;
if (document.cookie.indexOf(disableStr + '=true') > -1) {
  window[disableStr] = true;
}
function gaOptout() {
  document.cookie = disableStr + '=true; expires=Thu, 31 Dec 2099 23:59:59 UTC; path=/';
  window[disableStr] = true;
}
</script>
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-49679648-1', 'devabo.de');
  ga('set', 'anonymizeIp', true);
  ga('require', 'displayfeatures');
  ga('require', 'linkid', 'linkid.js');
  ga('send', 'pageview');

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
	return date.Format(time.RFC1123)
}

func (o *Output) Rss() string {
	date := o.DateNow()
	items := o.RssItems()
	return fmt.Sprintf(rss, date, items)
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
	<icon>
	</icon>
	<atom:link href="https://DevAbo.de/rss.xml" rel="self" type="application/rss+xml" />
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
    <guid isPermaLink="false">%s</guid>
    <description><![CDATA[%s]]></description>
    <content:encoded><![CDATA[%s]]></content:encoded>

    <media:thumbnail url="%s" />
    <media:content url="%s" medium="image">
      <media:title type="html">%s</media:title>
      <media:thumbnail url="%s" />
    </media:content>
  </item>
`
