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
