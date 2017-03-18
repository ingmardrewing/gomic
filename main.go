package main

import "github.com/ingmardrewing/gomic/comic"

func main() {

	comic := comic.NewComic()

	comic.AddPage("test", "http://wurst.de/")
	comic.AddPage("test1", "http://wurst.de/1")
	comic.AddPage("test2", "http://wurst.de/2")
	comic.AddPage("test3", "http://wurst.de/3")

	comic.ConnectPages()

	comic.PrintPages()
}

type Comic struct {
	name string
}
