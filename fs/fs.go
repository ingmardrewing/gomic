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

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
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
	o.writeJs()
	o.writeArchive()
	o.writeRss()
	o.writeAbout()
	o.writeImprint()
}

func (o *Output) writeAbout() {
	ah := NewDataHtml(about, config.Servedrootpath()+"/about.html")
	o.writeStringToFS(config.Rootpath()+"/about.html", ah.writePage("About"))
}

func (o *Output) writeImprint() {
	ah := NewDataHtml(imprint, config.Servedrootpath()+"/imprint.html")
	o.writeStringToFS(config.Rootpath()+"/imprint.html", ah.writePage("Imprint"))
}

func (o *Output) writeRss() {
	rss := newRss(o.comic)
	path := config.Rootpath() + "/feed/"
	filename := "rss.xml"
	o.prepareFileSystem(path)
	log.Println("Writing rss: ", path+filename)
	o.writeStringToFS(path+filename, rss.Rss())
}

func (o *Output) writeNarrativePages() {
	for _, p := range o.comic.GetPages() {
		o.writePageToFileSystem(p)
	}
}

func (o *Output) writeThumbnailFor(p *comic.Page) string {
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
	ah := NewDataHtml(arc, config.Servedrootpath()+"/archive.html")
	o.writeStringToFS(config.Rootpath()+"/archive.html", ah.writePage("Archive"))
}

func (o *Output) writeCss() {
	cg := newCss()
	p := config.Rootpath() + "/css"
	o.prepareFileSystem(p)
	fp := p + "/style.css"
	o.writeStringToFS(fp, cg.getCss())
}

func (o *Output) writeJs() {
	js := newJs()
	p := config.Rootpath() + "/js"
	o.prepareFileSystem(p)
	fp := p + "/script.js"
	o.writeStringToFS(fp, js.getJs())
}

func (o *Output) writePageToFileSystem(p *comic.Page) {
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
