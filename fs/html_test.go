package fs

import (
	"fmt"
	"testing"
)

func TestCreateTextNode(t *testing.T) {
	tn := createText("hello world")
	expected := "hello world"

	if tn.Render() != expected {
		t.Error(fmt.Sprintf(
			`text node doesn't contain %s, but %s `,
			expected,
			tn.Render()))
	}
}

func TestCreateNestedTextNode(t *testing.T) {
	n := createNode("p")
	n.AppendText("hello")

	txt := n.Render()
	expected := "<p>hello</p>"
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestCreateTwoNestedTextNodes(t *testing.T) {
	n := createNode("p")
	n.AppendText("hello")
	n.AppendText(" world")

	expected := "<p>hello world</p>"
	txt := n.Render()
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestCreateNestedHtmlNodes(t *testing.T) {
	div := createNode("div")
	div.AppendTag("p").AppendText("hello world")

	expected := "<div><p>hello world</p></div>"
	txt := div.Render()
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestAddAttributes(t *testing.T) {
	div := createNode("div")
	div.Attr("style", "background-color:black;").Attr("id", "test")

	expected := `<div style="background-color:black;" id="test"></div>`
	txt := div.Render()
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestCreateSimpleDom(t *testing.T) {
	html := createNode("html")
	html.Attr("lang", "en").AppendTag("head").AppendTag("title").AppendText("Hello World")
	html.AppendTag("body").AppendTag("div").AppendTag("p").AppendText("works!")

	expected := `<html lang="en"><head><title>Hello World</title></head><body><div><p>works!</p></div></body></html>`
	txt := html.Render()
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestCreateDomWithDoctype(t *testing.T) {
	h := newHtmlDoc()
	h.AddToHead(createNode("title").AppendText("Hello World"))
	h.AddToBody(createNode("div").AppendText("yo"))

	expected := `<!doctype html>
<html lang="en"><head><title>Hello World</title></head><body><div>yo</div></body></html>`

	txt := h.Render()
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestGenerateIconLinks(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addAndroidIconLinks()
	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><link rel="icon" type="image/png" sizes="192x192" href="/android-icon-192x192.png"></head><body></body></html>`
	if txt != expected {
		t.Error(fe(expected, txt))
	}

	hdw.addFaviconLinks()
	txt = hdw.Render()
	expected = `<!doctype html>
<html lang="en"><head><link rel="icon" type="image/png" sizes="192x192" href="/android-icon-192x192.png"><link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png"><link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png"><link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png"></head><body></body></html>`
	if txt != expected {
		t.Error(fe(expected, txt))
	}

	hdw.addAppleIconLinks()
	txt = hdw.Render()
	expected = `<!doctype html>
<html lang="en"><head><link rel="icon" type="image/png" sizes="192x192" href="/android-icon-192x192.png"><link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png"><link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png"><link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png"><link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png"><link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png"><link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png"><link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png"><link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png"><link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png"><link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png"><link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png"><link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png"></head><body></body></html>`
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestGenerateTitle(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.AddTitle("Hello World")
	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><title>Hello World</title></head><body></body></html>`
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestAddGoogleApiLinkToJQuery(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addGoogleApiLinkToJQuery()

	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script></head><body></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}

}

func TestAddCopyrightNotifier(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addCopyrightNotifier("2017")

	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head></head><body><div class="copyright">All content including but not limited to the art, characters, story, website design & graphics are &copy; copyright 2013-2017 Ingmar Drewing unless otherwise stated. All rights reserved. Do not copy, alter or reuse without expressed written permission.</div></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestAddCokieLawInfo(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addCookieLawInfo()

	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head></head><body><div id="cookie-law-info-bar">This website uses cookies to improve your experience. We'll assume you're ok with this, but you can opt-out if you wish.<a href="#" id="cookie_action_close_header" class="medium cli-plugin-button cli-plugin-main-button">Accept</a> <a href="http://www.drewing.de/blog/impressum-imprint/" id="CONSTANT_OPEN_URL" target="_blank" class="cli-plugin-main-link">Read More</a></div></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestAddFooterNavi(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.AddFooterNavi(`<ul><li><a href="test.html">test</a></li></ul>`)

	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head></head><body><footer><nav><ul><li><a href="test.html">test</a></li></ul></nav></footer></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestGoogleMetaData(t *testing.T) {
	google_data := []string{
		"name", "The Name or Title Here",
		"description", "This is the page description",
		"image", "http,//www.example.com/image.jpg",
	}
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addNameValueMetas(google_data)
	txt := hdw.Render()

	expected := `<!doctype html>
<html lang="en"><head><meta name="The Name or Title Here"><meta description="This is the page description"><meta image="http,//www.example.com/image.jpg"></head><body></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestOpenGraph(t *testing.T) {
	og_data := []string{
		"og,title", "Title Here",
		"og,type", "article",
		"og,url", "http,//www.example.com/",
		"og,image", "http,//example.com/image.jpg",
		"og,description", "Description Here",
		"og,site_name", "Site Name, i.e. Moz",
		"article,published_time", "2013-09-17T05,59,00+01,00",
		"article,modified_time", "2013-09-16T19,08,47+01,00",
		"article,section", "Article Section",
		"article,tag", "Article Tag",
		"fb,admins", "Facebook numberic ID",
	}

	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addNameValueMetas(og_data)
	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><meta og,title="Title Here"><meta og,type="article"><meta og,url="http,//www.example.com/"><meta og,image="http,//example.com/image.jpg"><meta og,description="Description Here"><meta og,site_name="Site Name, i.e. Moz"><meta article,published_time="2013-09-17T05,59,00+01,00"><meta article,modified_time="2013-09-16T19,08,47+01,00"><meta article,section="Article Section"><meta article,tag="Article Tag"><meta fb,admins="Facebook numberic ID"></head><body></body></html>`

	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func TestStandardMeta(t *testing.T) {
	hdw := newHtmlDocWrapper().(*htmlDocWrapper)
	hdw.addStandardMeta()
	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><meta viewport="width=device-width, initial-scale=1.0"><meta robots="index,follow"><meta author="Ingmar Drewing"><meta publisher="Ingmar Drewing"><meta keywords="web comic, comic, cartoon, sci fi, satire, parody, science fiction, action, software industry, pulp, nerd, geek"><meta DC.Subject="web comic, comic, cartoon, sci fi, science fiction, satire, parody action, software industry"><meta page-topic="Science Fiction Web-Comic"><meta http-equiv="content-type" content="text/html;charset=UTF-8"></head><body></body></html>`
	if txt != expected {
		t.Error(fe(expected, txt))
	}
}

func fe(expected string, actual string) string {
	return fmt.Sprintf(
		`didn't find %s, but %s `,
		expected,
		actual)
}
