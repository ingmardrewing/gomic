package comic

import (
	"testing"

	"github.com/ingmardrewing/gomic/config"
)

func TestIsRelevant(t *testing.T) {
	c := newComic()

	result := c.isRelevant(".DS_Store")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = c.isRelevant("thumb_DevAbode_0085.png")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = c.isRelevant("DevAbode_0085.png")
	if result {
		t.Error("Expected result to be true, but it is false")
	}
}

func TestIsNewFile(t *testing.T) {
	c := newComic()

	result := c.IsNewFile(".DS_Store")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = c.IsNewFile("DevAbode_0001.png")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = c.IsNewFile("DevAbode_0011.png")
	if result != true {
		t.Error("Expected result to be true, but it is false")
	}
}

func newComic() Comic {
	config.Read("/Users/drewing/Sites/gomic.yaml")
	pages := createPages()
	return Comic{config.Rootpath(), pages}
}

func createPages() []*Page {
	return []*Page{
		NewPage(
			"#1 A Step in the dark",
			"/2013/08/01/a-step-in-the-dark",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0001.png",
			"8 http://devabo.de/?p=8",
			"Act I"),
		NewPage(
			"#2 Negotiation",
			"/2013/08/31/negotiation",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0002.png",
			"35 http://devabo.de/?p=35",
			"Act I"),
		NewPage(
			"#3 Weapon of choice",
			"/2013/08/31/weapon-of-choice",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0003.png",
			"79 http://devabo.de/?p=79",
			"Act I"),
		NewPage(
			"#4 Cassandra",
			"/2013/10/15/cassandra",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0004.png",
			"158 http://devabo.de/?p=158",
			"Act I"),
		NewPage(
			"#5 Super saver",
			"/2013/11/01/super-saver",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0005.png",
			"162 http://devabo.de/?p=16",
			"Act I"),
		NewPage(
			"#6 Home, sweet home",
			"/2013/11/15/home-sweet-home",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0006.png",
			"177 http://devabo.de/?p=177",
			"Act I"),
		NewPage(
			"#7 The high council",
			"/2013/12/01/the-high-council",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0007.png",
			"192 http://devabo.de/?p=192",
			"Act I"),
		NewPage(
			"#8 Secrecy",
			"/2013/12/15/secrecy",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0008.png",
			"200 http://devabo.de/?p=200",
			"Act I"),
		NewPage(
			"#9 Curiosity",
			"/2014/01/01/curiosity",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0009.png",
			"208 http://devabo.de/?p=208",
			"Act I"),
		NewPage(
			"#10 Welcome to the machine",
			"/2014/01/15/welcome-to-the-machine",
			"https://devabode-us.s3.amazonaws.com/comicstrips/DevAbode_0010.png",
			"214 http://devabo.de/?p=214",
			"Act I")}
}
