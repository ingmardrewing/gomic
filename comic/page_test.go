package comic

import "testing"

func TestGetPath(t *testing.T) {
	expected := "/2017/04/29/#85-test"
	title := "#85-test"
	y := 2017
	m := 04
	d := 29
	actual := getPath(title,y,m,d)
	if actual != expected{
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestDisqusId(t *testing.T) {
	expected := "20170429 https://DevAbo.de/?p=20170429"
	y := 2017
	m := 04
	d := 29
	actual := getDisqusId(y,m,d)
	if actual != expected{
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestCreatePathTitleFromTitle(t *testing.T) {
	expected := "#85-The-Test"
	title := "#85 The  Test $"
	actual := createPathTitleFromTitle(title)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}
