package fs

import (
	"fmt"
	"testing"
)

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
	hdw.addTitle("Hello World")
	txt := hdw.Render()
	expected := `<!doctype html>
<html lang="en"><head><title>Hello World</title></head><body></body></html>`
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
