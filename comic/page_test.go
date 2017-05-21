package comic

import "testing"

func TestGetPath(t *testing.T) {
	expected := "/2017/04/29/85-test"
	title := "#85-test"
	y := 2017
	m := 04
	d := 29
	actual := getPath(title, y, m, d)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestDisqusId(t *testing.T) {
	expected := "20170429 https://DevAbo.de/?p=20170429"
	y := 2017
	m := 04
	d := 29
	actual := getDisqusId(y, m, d)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestCreatePathTitleFromTitle(t *testing.T) {
	title := "#85 The  Test $"
	expected := "85-The-Test"
	actual := createPathTitleFromTitle(title)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestImageFilename(t *testing.T) {
	expected := "DevAbode_0085.png"
	p := getPage()
	actual := p.ImageFilename()
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestProdUrl(t *testing.T) {
	expected := "https://devabo.de/2017/04/19/85-Test"
	p := getPage()
	actual := p.ProdUrl()
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestTitle(t *testing.T) {
	expected := "#85-Test"
	p := getPage()
	actual := p.Title()
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestPageDisqusId(t *testing.T) {
	expected := "20170419 http://DevAbo.de/?p=20170429"
	p := getPage()
	actual := p.DisqusId()
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func getPage() *Page {
	return NewPage("#85-Test", "/2017/04/19/85-Test", "http://localhost/DevAbode_0085.png", "20170419 http://DevAbo.de/?p=20170429", "III")
}
