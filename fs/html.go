package fs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
)

type node interface {
	Append(n node) node
	AppendText(txt string) node
	AppendTag(txt string) node
	Attr(name string, value string) node
	Render() string
	RenderReadable() string
	renderReadable(indent string) string
}

type htmlNode struct {
	name       string
	txt        string
	children   []node
	attributes []string
}

func createText(txt string) node {
	return &htmlNode{"", txt, []node{}, []string{}}
}

func createNode(name string) node {
	return &htmlNode{name, "", []node{}, []string{}}
}

func (n *htmlNode) setTxt(txt string) node {
	n.txt = txt
	return n
}

func (n *htmlNode) Render() string {
	attrs := n.getAttrs()
	if len(n.txt) > 0 {
		return n.txt
	} else if len(n.children) == 0 && isStandAloneTag(n.name) {
		return fmt.Sprintf("<%s%s>", n.name, attrs)
	}
	txt := n.getInner()
	return fmt.Sprintf("<%s%s>%s</%s>", n.name, attrs, txt, n.name)
}

func (n *htmlNode) RenderReadable() string {
	if len(n.txt) > 0 {
		return n.txt
	}
	txt := n.getInnerReadable("  ")
	attrs := n.getAttrs()
	return fmt.Sprintf("<%s%s>%s\n</%s>",
		n.name, attrs,
		txt,
		n.name)
}

func (n *htmlNode) renderReadable(indent string) string {
	attrs := n.getAttrs()
	if len(n.txt) > 0 {
		return "\n" + indent + n.txt
	} else if len(n.children) == 0 && isStandAloneTag(n.name) {
		return fmt.Sprintf("\n%s<%s%s />",
			indent, n.name, attrs)
	}
	txt := n.getInnerReadable(indent + "  ")
	return fmt.Sprintf("\n%s<%s%s>%s\n%s</%s>",
		indent, n.name, attrs,
		txt,
		indent, n.name)
}

func (n *htmlNode) getInnerReadable(indent string) string {
	txt := ""
	for _, c := range n.children {
		txt += indent + c.renderReadable(indent)
	}
	return txt
}

func (n *htmlNode) getInner() string {
	txt := ""
	for _, c := range n.children {
		txt += c.Render()
	}
	return txt
}

func (n *htmlNode) getAttrs() string {
	if len(n.attributes) > 0 {
		return " " + strings.Join(n.attributes, " ")
	}
	return ""
}

func (n *htmlNode) Append(nd node) node {
	n.children = append(n.children, nd)
	return nd
}

func (n *htmlNode) Attr(name string, value string) node {
	n.attributes = append(
		n.attributes,
		fmt.Sprintf(`%s="%s"`, name, value))
	return n
}

func (n *htmlNode) AppendText(txt string) node {
	n.children = append(n.children, createText(txt))
	return n
}

func (n *htmlNode) AppendTag(nn string) node {
	cn := createNode(nn)
	n.Append(cn)
	return cn
}

type htmlDoc interface {
	Render() string
	AddToHead(n node)
	AddToBody(n node)
}

type html struct {
	doctype string
	head    node
	body    node
}

func newHtmlDoc() htmlDoc {
	return &html{
		"<!doctype html>",
		createNode("head"),
		createNode("body")}
}

func (h *html) Render() string {
	root := createNode("html").Attr("lang", "en")
	root.Append(h.head)
	root.Append(h.body)
	return h.doctype + "\n" + root.Render()
}

func (h *html) AddToHead(n node) {
	h.head.Append(n)
}

func (h *html) AddToBody(n node) {
	h.body.Append(n)
}

func isStandAloneTag(tagname string) bool {
	standalones := []string{"img", "link", "meta"}
	for _, t := range standalones {
		if t == tagname {
			return true
		}
	}
	return false
}

type htmlDocWrapperI interface {
	Render() string
	Version() string
	AddToHead(n node)
	AddToBody(n node)
	AddTitle(txt string)
	AddCopyrightNotifier(year string)
	AddFooterNavi(txt string)
	AddNameValueMetas(mataData []string)
	AddCookieLawInfo()
	AddNewsletter()
	Init()
}

type htmlDocWrapper struct {
	htmlDoc htmlDoc
}

func newHtmlDocWrapper() htmlDocWrapperI {
	return &htmlDocWrapper{newHtmlDoc()}
}

func (hdw *htmlDocWrapper) Init() {
	hdw.addStandardMeta()
	hdw.addAndroidIconLinks()
	hdw.addFaviconLinks()
	hdw.addAppleIconLinks()
	hdw.addGoogleApiLinkToJQuery()
}

func (hdw *htmlDocWrapper) addGoogleApiLinkToJQuery() {
	hdw.htmlDoc.AddToHead(createNode("script").Attr("src", "https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"))
}

func (hdw *htmlDocWrapper) AddTitle(txt string) {
	hdw.htmlDoc.AddToHead(createNode("title").AppendText(txt))
}

func (hdw *htmlDocWrapper) addStandardMeta() {
	name_metas := []string{
		"viewport", "width=device-width, initial-scale=1.0",
		"robots", "index,follow",
		"author", "Ingmar Drewing",
		"publisher", "Ingmar Drewing",
		"keywords", "web comic, comic, cartoon, sci fi, satire, parody, science fiction, action, software industry, pulp, nerd, geek",
		"DC.Subject", "web comic, comic, cartoon, sci fi, science fiction, satire, parody action, software industry",
		"page-topic", "Science Fiction Web-Comic",
	}
	hdw.AddNameValueMetas(name_metas)
	hdw.htmlDoc.AddToHead(createNode("meta").Attr("http-equiv", "content-type").Attr("content", "text/html;charset=UTF-8"))
}

func (hdw *htmlDocWrapper) AddNameValueMetas(metaData []string) {
	for i := 0; i < len(metaData); i += 2 {
		m := createNode("meta")
		m.Attr(metaData[i], metaData[i+1])
		hdw.htmlDoc.AddToHead(m)
	}
}

func (hdw *htmlDocWrapper) AddCopyrightNotifier(year string) {
	hdw.htmlDoc.AddToBody(createNode("div").Attr("class", "copyright").AppendText(`All content including but not limited to the art, characters, story, website design & graphics are &copy; copyright 2013-` + year + ` Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.`))
}

func (hdw *htmlDocWrapper) AddCookieLawInfo() {
	hdw.htmlDoc.AddToBody(createNode("div").Attr("id", "cookie-law-info-bar").AppendText(`This website uses cookies to improve your experience. We'll assume you're ok with this, but you can opt-out if you wish.<a href="#" id="cookie_action_close_header" class="medium cli-plugin-button cli-plugin-main-button">Accept</a> <a href="http://www.drewing.de/blog/impressum-imprint/" id="CONSTANT_OPEN_URL" target="_blank" class="cli-plugin-main-link">Read More</a>`))
}

func (hdw *htmlDocWrapper) AddFooterNavi(navi string) {
	n := createNode("footer")
	n.AppendTag("nav").AppendText(navi)
	hdw.htmlDoc.AddToBody(n)
}

func (hdw *htmlDocWrapper) AddNewsletter() {
	n := createNode("div")
	n.Attr("class", "nl_container nl_container_hidden")
	hdw.htmlDoc.AddToBody(n)
}

func (hdw *htmlDocWrapper) addFaviconLinks() {
	iconSizes := []string{
		"32x32",
		"96x96",
		"16x16",
	}
	for _, s := range iconSizes {
		l := createNode("link")
		l.Attr("rel", "icon")
		l.Attr("type", "image/png")
		l.Attr("sizes", s)
		l.Attr("href", "/icons/favicon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) Version() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02dT%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (hdw *htmlDocWrapper) AddToBody(n node) {
	hdw.htmlDoc.AddToBody(n)
}

func (hdw *htmlDocWrapper) AddToHead(n node) {
	hdw.htmlDoc.AddToHead(n)
}

func (hdw *htmlDocWrapper) addAndroidIconLinks() {
	androidSizes := []string{
		"192x192",
	}
	for _, s := range androidSizes {
		l := createNode("link")
		l.Attr("rel", "icon")
		l.Attr("type", "image/png")
		l.Attr("sizes", s)
		l.Attr("href", "/icons/android-icon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) addAppleIconLinks() {
	appleSizes := []string{
		"57x57",
		"60x60",
		"72x72",
		"76x76",
		"114x114",
		"120x120",
		"144x144",
		"152x152",
		"180x180",
	}
	for _, s := range appleSizes {
		l := createNode("link")
		l.Attr("rel", "apple-touch-icon")
		l.Attr("sizes", s)
		l.Attr("href", "/icons/apple-icon-"+s+".png")
		hdw.htmlDoc.AddToHead(l)
	}
}

func (hdw *htmlDocWrapper) Render() string {
	return hdw.htmlDoc.Render()
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

	js := newJs()
	h += js.getAnalytics()
	return h
}

func (html *HTML) getCssLink() string {
	path := config.Servedrootpath() + "/css/style.css?version=" + html.version()
	format := `<link rel="stylesheet" href="%s" type="text/css">`
	return fmt.Sprintf(format, path)
}

func DateNow() string {
	date := time.Now()
	return date.Format(time.RFC1123Z)
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
	hdw.AddCookieLawInfo()

	hdw.AddFooterNavi(ah.getFooterNavi())
	hdw.AddNewsletter()

	return hdw.Render()
}

type NarrativePageHtml struct {
	HTML
	p          *comic.Page
	title      string
	meta       string
	csslink    string
	img        string
	navi       string
	footerNavi string
}

func NewNarrativePageHtml(p *comic.Page) *NarrativePageHtml {
	return &NarrativePageHtml{HTML{}, p, "", "", "", "", "", ""}
}

func (h *NarrativePageHtml) writePage() string {

	hdw := newHtmlDocWrapper()
	hdw.Init()

	css_path := config.Servedrootpath() + "/css/style.css?version=" + hdw.Version()
	hdw.AddToHead(createNode("link").Attr("rel", "stylesheet").Attr("href", css_path).Attr("type", "text/css"))
	hdw.AddToHead(createNode("link").Attr("rel", "canonical").Attr("href", h.p.GetPath()))

	js_path := config.Servedrootpath() + "/js/script.js?version=" + hdw.Version()
	hdw.AddToHead(createNode("script").Attr("src", js_path).Attr("type", "text/javascript").Attr("language", "javascript"))
	hdw.AddTitle("DevAbo.de | Graphic Novel | " + h.p.GetTitle())

	header := createNode("header").AppendText(h.getHeaderHtml())
	hdw.AddToBody(header)

	main := createNode("main")
	main.AppendText(h.getContent())
	main.AppendText(h.getNaviHtml())
	main.AppendText(h.getDisqus())
	hdw.AddToBody(main)

	hdw.AddCopyrightNotifier(strconv.Itoa(time.Now().Year()))
	hdw.AddCookieLawInfo()

	hdw.AddFooterNavi(h.getFooterNavi())
	hdw.AddNewsletter()

	desc := h.p.GetDescription()
	if len(desc) == 0 {
		desc = "A dystopian sci-fi webcomic about the life of software developers"
	}
	hdw.AddToHead(createNode("meta").Attr("property", "og:title").Attr("content", h.p.GetTitle()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:url").Attr("content", h.p.GetPath()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:image").Attr("content", h.p.GetImgUrl()))
	hdw.AddToHead(createNode("meta").Attr("property", "og:description").Attr("content", desc))
	hdw.AddToHead(createNode("meta").Attr("property", "og:site_name").Attr("content", "DevAbo.de"))
	hdw.AddToHead(createNode("meta").Attr("property", "og:type").Attr("content", "article"))
	hdw.AddToHead(createNode("meta").Attr("property", "article:published_time").Attr("content", h.p.GetDateFromFSPath()))
	hdw.AddToHead(createNode("meta").Attr("property", "article:modified_time").Attr("content", h.p.GetDateFromFSPath()))
	hdw.AddToHead(createNode("meta").Attr("property", "article:section").Attr("content", "Science-Fiction"))
	hdw.AddToHead(createNode("meta").Attr("property", "article:tag").Attr("content", "comic, graphic novel, webcomic, science-fiction, sci-fi"))

	hdw.AddToHead(createNode("meta").Attr("itemprop", "name").Attr("content", h.p.GetTitle()))
	hdw.AddToHead(createNode("meta").Attr("itemprop", "name").Attr("description", desc))
	hdw.AddToHead(createNode("meta").Attr("itemprop", "image").Attr("content", h.p.GetImgUrl()))

	hdw.AddToHead(createNode("meta").Attr("name", "twitter:card").Attr("content", "summary_large_image"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:site").Attr("content", "@devabo_de"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:title").Attr("content", h.p.GetTitle()))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:text:description").Attr("content", desc))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:creator").Attr("content", "@ingmardrewing"))
	hdw.AddToHead(createNode("meta").Attr("name", "twitter:image").Attr("content", h.p.GetImgUrl()))

	return hdw.Render()
}

func (h *NarrativePageHtml) getContent() string {
	f := `<img src="%s" width="800" height="1334" alt="">`
	html := fmt.Sprintf(f, h.p.GetImgUrl())
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
	hl := h.getHeadline(h.p.GetTitle())
	s := config.Servedrootpath()
	return fmt.Sprintf(`
	<a href="%s" class="home"><!--DevAbo.de--></a>
    <a href="%s/2013/08/01/a-step-in-the-dark/" class="orange">New Reader? Start here!</a>
	%s`, s, s, hl)
}

func (h *NarrativePageHtml) getDisqus() string {
	js := newJs()
	return js.getDisqus(
		h.p.GetTitle(),
		h.getDisqusUrl(),
		h.getDisqusIdentifier(),
	)
}

func (h *NarrativePageHtml) getDisqusIdentifier() string {
	if len(h.p.GetDisqusIdentifier()) > 0 {
		return h.p.GetDisqusIdentifier()
	}
	return h.p.GetPath()
}

func (h *NarrativePageHtml) getDisqusUrl() string {
	return h.p.GetPath() + "/"
}

var imprint = `
<h3>Angaben nach TDG:</h3>

<p>Dieses Impressum gilt für die Website devabo.de, so wie z.B. die zugehörige Facebook-Page <a href="https://www.facebook.com/devabo.de">https://www.facebook.com/devabo.de</a> und das Facebook-Profil <a href="https://www.facebook.com/ingmar.drewing">https://www.facebook.com/ingmar.drewing</a> sowie die Google-Plus Seite <a href="https://plus.google.com/107755781256885973202/posts">https://plus.google.com/107755781256885973202/posts</a> der Twitter-Account unter <a href="https://twitter.com/ingmardrewing">https://twitter.com/ingmardrewing</a> und alle weiteren Profile und Websites von Ingmar Drewing, wie auch www.devabo.de .</p>

<h4>Redaktionell verantwortlich ist:</h4>
<p>
Ingmar Drewing<br>
(Dipl. Kommunkationsdesigner /FH /BRD)<br>
Schulberg 8<br>
65183 Wiesbaden<br>
</p>

<p>
Telefon: 0173-3076520<br>
E-Mail: ingmar-at-drewing-punkt-de
<p>


<h3>Haftungsausschluss</h3>

<h4>1. Inhalt des Onlineangebotes</h4>
<p>Der Autor übernimmt keinerlei Gewähr für die Aktualität, Korrektheit, Vollständigkeit oder Qualität der bereitgestellten Informationen. Haftungsansprüche gegen den Autor, welche sich auf Schäden materieller oder ideeller Art beziehen, die durch die Nutzung oder Nichtnutzung der dargebotenen Informationen bzw. durch die Nutzung fehlerhafter und unvollständiger Informationen verursacht wurden, sind grundsätzlich ausgeschlossen, sofern seitens des Autors kein nachweislich vorsätzliches oder grob fahrlässiges Verschulden vorliegt.</p>
<p>Alle Angebote sind freibleibend und unverbindlich. Der Autor behält es sich ausdrücklich vor, Teile der Seiten oder das gesamte Angebot ohne gesonderte Ankündigung zu verändern, zu ergänzen, zu löschen oder die Veröffentlichung zeitweise oder endgültig einzustellen.</p>

<h4>2. Verweise und Links</h4>
<p>Bei direkten oder indirekten Verweisen auf fremde Webseiten ("Hyperlinks"), die außerhalb des Verantwortungsbereiches des Autors liegen, würde eine Haftungsverpflichtung ausschließlich in dem Fall in Kraft treten, in dem der Autor von den Inhalten Kenntnis hat und es ihm technisch möglich und zumutbar wäre, die Nutzung im Falle rechtswidriger Inhalte zu verhindern.</p>
<p>Der Autor erklärt hiermit ausdrücklich, dass zum Zeitpunkt der Linksetzung keine illegalen Inhalte auf den zu verlinkenden Seiten erkennbar waren. Auf die aktuelle und zukünftige Gestaltung, die Inhalte oder die Urheberschaft der verlinkten/verknüpften Seiten hat der Autor keinerlei Einfluss. Deshalb distanziert er sich hiermit ausdrücklich von allen Inhalten aller verlinkten /verknüpften Seiten, die nach der Linksetzung verändert wurden. Diese Feststellung gilt für alle innerhalb des eigenen Internetangebotes gesetzten Links und Verweise sowie für Fremdeinträge in vom Autor eingerichteten Gästebüchern, Diskussionsforen, Linkverzeichnissen, Mailinglisten und in allen anderen Formen von Datenbanken, auf deren Inhalt externe Schreibzugriffe möglich sind. Für illegale, fehlerhafte oder unvollständige Inhalte und insbesondere für Schäden, die aus der Nutzung oder Nichtnutzung solcherart dargebotener Informationen entstehen, haftet allein der Anbieter der Seite, auf welche verwiesen wurde, nicht derjenige, der über Links auf die jeweilige Veröffentlichung lediglich verweist.</p>

<h4>3. Urheber- und Kennzeichenrecht</h4>
<p>Der Autor ist bestrebt, in allen Publikationen die Urheberrechte der verwendeten Bilder, Grafiken, Tondokumente, Videosequenzen und Texte zu beachten, von ihm selbst erstellte Bilder, Grafiken, Tondokumente, Videosequenzen und Texte zu nutzen oder auf lizenzfreie Grafiken, Tondokumente, Videosequenzen und Texte zurückzugreifen.</p>
<p>Alle innerhalb des Internetangebotes genannten und ggf. durch Dritte geschützten Marken- und Warenzeichen unterliegen uneingeschränkt den Bestimmungen des jeweils gültigen Kennzeichenrechts und den Besitzrechten der jeweiligen eingetragenen Eigentümer. Allein aufgrund der bloßen Nennung ist nicht der Schluss zu ziehen, dass Markenzeichen nicht durch Rechte Dritter geschützt sind!</p>
<p>Das Copyright für veröffentlichte, vom Autor selbst erstellte Objekte bleibt allein beim Autor der Seiten. Eine Vervielfältigung oder Verwendung solcher Grafiken, Tondokumente, Videosequenzen und Texte in anderen elektronischen oder gedruckten Publikationen ist ohne ausdrückliche Zustimmung des Autors nicht gestattet.</p>

<h4>4. Datenschutz</h4>
<p>Sofern innerhalb des Internetangebotes die Möglichkeit zur Eingabe persönlicher oder geschäftlicher Daten (Emailadressen, Namen, Anschriften) besteht, so erfolgt die Preisgabe dieser Daten seitens des Nutzers auf ausdrücklich freiwilliger Basis. Die Inanspruchnahme und Bezahlung aller angebotenen Dienste ist - soweit technisch möglich und zumutbar - auch ohne Angabe solcher Daten bzw. unter Angabe anonymisierter Daten oder eines Pseudonyms gestattet. Die Nutzung der im Rahmen des Impressums oder vergleichbarer Angaben veröffentlichten Kontaktdaten wie Postanschriften, Telefon- und Faxnummern sowie Emailadressen durch Dritte zur Übersendung von nicht ausdrücklich angeforderten Informationen ist nicht gestattet. Rechtliche Schritte gegen die Versender von sogenannten Spam-Mails bei Verstössen gegen dieses Verbot sind ausdrücklich vorbehalten.</p>

<h4>5. Rechtswirksamkeit dieses Haftungsausschlusses</h4>
<p>Dieser Haftungsausschluss ist als Teil des Internetangebotes zu betrachten, von dem aus auf diese Seite verwiesen wurde. Sofern Teile oder einzelne Formulierungen dieses Textes der geltenden Rechtslage nicht, nicht mehr oder nicht vollständig entsprechen sollten, bleiben die übrigen Teile des Dokumentes in ihrem Inhalt und ihrer Gültigkeit davon unberührt.</p>

<h4>6. Google Analytics (Text übernommen von <a href="http://www.datenschutzbeauftragter-info.de">www.datenschutzbeauftragter-info.de</a>)</h4>
<p>Diese Website benutzt Google Analytics, einen Webanalysedienst der Google Inc. („Google“). Google Analytics verwendet sog. „Cookies“, Textdateien, die auf Ihrem Computer gespeichert werden und die eine Analyse der Benutzung der Website durch Sie ermöglichen. Die durch den Cookie erzeugten Informationen über Ihre Benutzung dieser Website werden in der Regel an einen Server von Google in den USA übertragen und dort gespeichert. Im Falle der Aktivierung der IP-Anonymisierung auf dieser Website, wird Ihre IP-Adresse von Google jedoch innerhalb von Mitgliedstaaten der Europäischen Union oder in anderen Vertragsstaaten des Abkommens über den Europäischen Wirtschaftsraum zuvor gekürzt. Nur in Ausnahmefällen wird die volle IP-Adresse an einen Server von Google in den USA übertragen und dort gekürzt. Im Auftrag des Betreibers dieser Website wird Google diese Informationen benutzen, um Ihre Nutzung der Website auszuwerten, um Reports über die Websiteaktivitäten zusammenzustellen und um weitere mit der Websitenutzung und der Internetnutzung verbundene Dienstleistungen gegenüber dem Websitebetreiber zu erbringen. Die im Rahmen von Google Analytics von Ihrem Browser übermittelte IP-Adresse wird nicht mit anderen Daten von Google zusammengeführt. Sie können die Speicherung der Cookies durch eine entsprechende Einstellung Ihrer Browser-Software verhindern; wir weisen Sie jedoch darauf hin, dass Sie in diesem Fall gegebenenfalls nicht sämtliche Funktionen dieser Website vollumfänglich werden nutzen können. Sie können darüber hinaus die Erfassung der durch das Cookie erzeugten und auf Ihre Nutzung der Website bezogenen Daten (inkl. Ihrer IP-Adresse) an Google sowie die Verarbeitung dieser Daten durch Google verhindern, indem sie das unter dem folgenden Link (<a href="http://tools.google.com/dlpage/gaoptout?hl=de">http://tools.google.com/dlpage/gaoptout?hl=de</a>) verfügbare Browser-Plugin herunterladen und installieren.</p>

<p>Sie können die Erfassung durch Google Analytics verhindern, indem Sie auf folgenden Link klicken. Es wird ein Opt-Out-Cookie gesetzt, der die zukünftige Erfassung Ihrer Daten beim Besuch dieser Website verhindert:</p>
<a href="javascript:gaOptout()">Google Analytics deaktivieren</a>

<p>Nähere Informationen zu Nutzungsbedingungen und Datenschutz finden Sie unter <a href="http://www.google.com/analytics/terms/de.html">http://www.google.com/analytics/terms/de.html</a> bzw. unter <a href="http://www.google.com/intl/de/analytics/privacyoverview.html">http://www.google.com/intl/de/analytics/privacyoverview.html</a>. Wir weisen Sie darauf hin, dass auf dieser Website Google Analytics um den Code „gat._anonymizeIp();“ erweitert wurde, um eine anonymisierte Erfassung von IP-Adressen (sog. IP-Masking) zu gewährleisten.</p>


<h3>Disclaimer</h3>

<h4>1. Content</h4>
<p>The author reserves the right not to be responsible for the topicality, correctness, completeness or quality of the information provided. Liability claims regarding damage caused by the use of any information provided, including any kind of information which is incomplete or incorrect,will therefore be rejected.</p>
<p>All offers are not-binding and without obligation. Parts of the pages or the complete publication including all offers and information might be extended, changed or partly or completely deleted by the author without separate announcement.</p>

<h4>2. Referrals and links</h4>
<p>The author is not responsible for any contents linked or referred to from his pages - unless he has full knowledge of illegal contents and would be able to prevent the visitors of his site fromviewing those pages. If any damage occurs by the use of information presented there, only the author of the respective pages might be liable, not the one who has linked to these pages. Furthermore the author is not liable for any postings or messages published by users of discussion boards, guestbooks or mailinglists provided on his page.</p>

<h4>3. Copyright</h4>
<p>The author intended not to use any copyrighted material for the publication or, if not possible, to indicate the copyright of the respective object.</p>
<p>The copyright for any material created by the author is reserved. Any duplication or use of objects such as images, diagrams, sounds or texts in other electronic or printed publications is not permitted without the author's agreement.</p>

<h4>4. Privacy policy</h4>
<p>If the opportunity for the input of personal or business data (email addresses, name, addresses) is given, the input of these data takes place voluntarily. The use and payment of all offered services are permitted - if and so far technically possible and reasonable - without specification of any personal data or under specification of anonymized data or an alias. The use of published postal addresses, telephone or fax numbers and email addresses for marketing purposes is prohibited, offenders sending unwanted spam messages will be punished.</p>

<h4>5. Legal validity of this disclaimer</h4>
<p>This disclaimer is to be regarded as part of the internet publication which you were referred from. If sections or individual terms of this statement are not legal or correct, the content or validity of the other parts remain uninfluenced by this fact.</p>

<h4>6. Google Analytics (Text by <a href="http://www.datenschutzbeauftragter-info.de">www.datenschutzbeauftragter-info.de</a>)</h4>
<p>This website uses Google Analytics, a web analytics service provided by Google, Inc. (“Google”).  Google Analytics uses “cookies”, which are text files placed on your computer, to help the website analyze how users use the site. The information generated by the cookie about your use of the website (including your IP address) will be transmitted to and stored by Google on servers in the United States.  In case of activation of the IP anonymization, Google will truncate/anonymize the last octet of the IP address for Member States of the European Union as well as for other parties to the Agreement on the European Economic Area.  Only in exceptional cases, the full IP address is sent to and shortened by Google servers in the USA.  On behalf of the website provider Google will use this information for the purpose of evaluating your use of the website, compiling reports on website activity for website operators and providing other services relating to website activity and internet usage to the website provider.  Google will not associate your IP address with any other data held by Google.  You may refuse the use of cookies by selecting the appropriate settings on your browser. However, please note that if you do this, you may not be able to use the full functionality of this website.  Furthermore you can prevent Google’s collection and use of data (cookies and IP address) by downloading and installing the browser plug-in available under <a href="https://tools.google.com/dlpage/gaoptout?hl=en-GB">https://tools.google.com/dlpage/gaoptout?hl=en-GB</a>.</p>
<p>You can refuse the use of Google Analytics by clicking on the following link. An opt-out cookie will be set on the computer, which prevents the future collection of your data when visiting this website:</p>
<a href="javascript:gaOptout()">Disable Google Analytics</a>
<p>Further information concerning the terms and conditions of use and data privacy can be found at <a href="http://www.google.com/analytics/terms/gb.html">http://www.google.com/analytics/terms/gb.html</a> or at <a href="http://www.google.com/intl/en_uk/analytics/privacyoverview.html">http://www.google.com/intl/en_uk/analytics/privacyoverview.html</a>.  Please note that on this website, Google Analytics code is supplemented by “gat._anonymizeIp();” to ensure an anonymized collection of IP addresses (so called IP-masking).</p>
`

var about = `
devabo.de is a software developers fever dream, a nightmarish science-fiction comic by Ingmar Drewing. You can find the associated facebook-page at <a href="https://www.facebook.com/www.devabo.de">https://www.facebook.com/www.devabo.de</a> and the associated twitter account at <a href="twitter.com/#!/devabo_de">twitter.com/#!/devabo_de</a>.<br />
DevAbo.de will be updated every 1st and 15th day of every month, though I am trying right now to speed the production up to a weekly release cycle (but that's still beta).


<h3><a name="Bram">Bram</a></h3>
Bram has spent over a millennium in cryo stasis. <a href="#Ada">Ada</a> found him in the ancient ruins and ended his cryostatic slumber out of curiosity (and the possibility that he might have taken something of value into the cryo capsule). He woke up healthy, though a big part of his episodic memory is lost to him and resurfaces partially and slowly. <br /><br />
Some parts of his past that came back to him showed that he was some kind of <a href="http://devabo.de/2014/04/01/flashback/">technical officer</a>. He first didn't recall the aggressor he was fighting against. The memory of this came back to him while he was teaching Ada how to get in touch in with the calculating space and she <a href="http://devabo.de/2014/07/15/22-backup/">accidentally changed parts of his memory</a>.<br /><br />
Bram was about 35 years old when he was put into cryo slumber. He is still failing to remember the reason and circumstances of him being put into cryostasis, though it's likely that it has something to do with the war he was fighting in the past &hellip;<br /><br /><br /><br />


<h3><a name="Ada">Ada</a></h3>
Ada is a developer of the abode as well as an elite fighter. At the beginning of the story she is 27 years old.<br />
On routine checks along the outer defense perimeters of the abode she found entries to some rather well preserved buildings of the ancients and started to sell the artifacts she found there to <a href="#MasterBranch">Master Branch</a>. The business relationship developed and she regularly helped to retrieve artifacts for Master Branch. The business already brought her into conflict with her superiors and though she usually tries to keep out of trouble with the administration of the abode, she sneaked out into the ruins occasionally to "check for some strange client activity", as she put it in her report.<br />
Apart from this she takes her duty very seriously and is a good comrade. If a friend of her is in danger she's more than ready to risk her own life to free him. And she expects the same behaviour from everyone of her comrades.
<br /><br /><br /><br />

<h3><a name="MasterBranch">Master Branch</a></h3>
 Master Branch is a thrirty years old JMonk and, like all of these pious people, believes strongly in <a href="http://en.wikipedia.org/wiki/Type_system#Static_type-checking">static typing</a>. Since that faith is mercilessly tested every time reality interferes with their believe system, the JMonks are also on the lookout for a sign from a higher power. They have a prophecy that one day a man would come, a Messiah, who will bring them true productivity. But, until this comes true, their only joy will be the incredible beauty of generics and jverbosity&trade;.<br /><br />
However, they still managed to create a machine that emenates fields of unproductivity. Though the specs didn't say the machine would do this, many of the JMonks hope that the machine might perhaps prove useful after all, one day.
Because of the difficulties mentioned above, Master Branch made a deal with a developer, <a href="#Ada">Ada</a>, whom he bid to go and search the ancient ruins for useful artifacts. He hoped to reverse engineer the artifacts and maybe find a way to become productive. Ada found and delivered several artifacts, which seemed interesting. Unfortunately they didn't reveal their usefulness yet. <br />
The most peculiar thing she found in the ruins she didn't deliver to the monks at all ...<br /><br /><br /><br />

<h3>Clients</h3>
Clients are ruthless and aggressive and also rather dim. That's making them less of a threat, as long as you are sufficiently quick witted. Nevertheless a greater pack of them can be quite distressing for a developer or consultant.
<br/><br/><br/><br/>


<h2>The Author</h2>

If you are interested in the other things I am drawing and writing you'll find some fragments on my <a href="http://www.drewing.de/blog">blog</a>. A word of warning: this blog contains material some people might consider nsfw.

`

const imageWrapperFormat = `<a href="%s" rel="next" title="%s">%s</a>`
const navWrapperFormat = `<nav>%s</nav>`
