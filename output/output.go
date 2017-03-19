package output

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

type Output struct {
	comic  *comic.Comic
	config *config.Config
}

func NewOutput(comic *comic.Comic, config *config.Config) *Output {
	return &Output{comic, config}
}

func (o *Output) WriteToFilesystem() {
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
	o.writeCss()
}

func (o *Output) writeCss() {
	p := o.config.Rootpath + "/css"
	o.prepareFileSystem(p)
	fp := p + "/style.css"
	o.writeStringToFS(fp, css)
}

func (o *Output) writePageToFileSystem(p *page.Page) {
	absPath := o.config.Rootpath + p.FSPath()
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
	path := h.p.ServedRootPath() + "/css/style.css?version=" + h.version()
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

	return fmt.Sprintf(htmlFormat, title, meta, css, header, content, navi, footerNavi)
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
<header>
	<a href="https://DevAbo.de/" class="home"><!--DevAbo.de--></a>
    <a href="https://devabo.de/2013/08/01/a-step-in-the-dark/" class="orange">New Reader? Start here!</a>
	%s
</header>, `, hl)
}

const css = `
header {
	width: 800px;
	margin: 0 auto;
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
const htmlFormat = `<!doctype html>
<html>
	<head>
		<title>%s</title>
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

<div id="disqus_thread">
     <div id="dsq-content">
         <ul id="dsq-comments">
            <li class="post pingback"><p>Pingback: <a href='https://DevAbo.de/2017/02/27/83-professionals/' rel='external nofollow' class='url'>DevAbo.de | Sci-Fi Webcomic and Graphic Novel | #83 Professionals &laquo;</a>)</p></li><!-- #comment-## -->
		</ul>
  </div>
</div>

<script type="text/javascript">
var disqus_url = 'https://DevAbo.de/2017/03/18/84-time-crystals/';
var disqus_identifier = '1235 https://DevAbo.de/?p=1235';
var disqus_container_id = 'disqus_thread';
var disqus_shortname = 'devabode';
var disqus_title = "#84 Time Crystals";
var disqus_config_custom = window.disqus_config;
var disqus_config = function () {
    /*
    All currently supported events:
    onReady: fires when everything is ready,
    onNewComment: fires when a new comment is posted,
    onIdentify: fires when user is authenticated
    */
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
    var dsq = document.createElement('script'); dsq.type = 'text/javascript';
    dsq.async = true;
    dsq.src = '//' + disqus_shortname + '.disqus.com/embed.js';
    (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
})();
</script>

<div class="spacer"></div>
		<footer><nav>%s</nav></footer>
	</body>
</html>
`
