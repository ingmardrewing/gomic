package page

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/ingmardrewing/gomic/config"
)

type Page struct {
	title, path, imgUrl, disqusId string
	first, prev, next, last       *Page
	meta, navi                    [][]string
}

func NewPageFromFilename(filename string) *Page {

	var title, path, imgUrl, disqusId string
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter title for %s: ", filename)
		title, _ = reader.ReadString('\n')
		title = strings.TrimSpace(title)

		whitespace := regexp.MustCompile(`\s+`)
		forbidden := regexp.MustCompile(`[^-A-Za-z0-9]`)
		trailingdash := regexp.MustCompile(`-$`)
		pathTitle := whitespace.ReplaceAllString(title, "-")
		pathTitle = forbidden.ReplaceAllString(pathTitle, "")
		pathTitle = trailingdash.ReplaceAllString(pathTitle, "")

		t := time.Now()
		y := t.Year()
		m := int(t.Month())
		d := t.Day()
		path = fmt.Sprintf("/%d/%02d/%02d/%s", y, m, d, pathTitle)

		id := y*10000 + m*100 + d
		disqusId = fmt.Sprintf("%d https://DevAbo.de/?p=%d", id, id)

		imgUrl = fmt.Sprintf("https://s3-us-west-1.amazonaws.com/devabode-us/comicstrips/%s", filename)

		summary := fmt.Sprintf("\ntitle: %s\npath: %s\ndisqusId: %s\nimgUrl: %s\n", title, path, disqusId, imgUrl)

		answer := AskUser(
			fmt.Sprintf(
				"Creating the following page:\n%s\nok? [yN]", summary))

		if answer {
			break
		}
		fmt.Println("Okay, let's try again ...")
	}
	// Initialize a session reading env vars
	sess, err := session.NewSession()

	// Create S3 service client
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	answer := AskUser("Proceed? [yN]")
	if answer {
		fmt.Println("continuing")
	}

	bucket := config.AwsBucket()
	localPathToFile := fmt.Sprintf("%s/%s", config.PngDir(), filename)
	remotePathToFile := fmt.Sprintf("%s/%s", config.AwsDir(), filename)
	file, err := os.Open(localPathToFile)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePathToFile),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	answer = AskUser("Proceed? [yN]")
	if answer {
		fmt.Println("continuing")
	}

	return &Page{title, path, imgUrl, disqusId,
		nil, nil, nil, nil, [][]string{}, [][]string{}}
}

func AskUser(question string) bool {
	fmt.Println(question)
	reader := bufio.NewReader(os.Stdin)
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)
	return confirmation == "y" || confirmation == "Y"
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func NewPage(
	title string,
	path string,
	imgUrl string,
	disqusId string) *Page {
	return &Page{title, path, imgUrl, disqusId,
		nil, nil, nil, nil, [][]string{}, [][]string{}}

}

func (p *Page) Filename() string {
	pathParts := strings.Split(p.imgUrl, "/")
	return pathParts[len(pathParts)-1]
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) DisqusIdentifier() string {
	return p.disqusId
}

func (p *Page) SetRels(first *Page, prev *Page, next *Page, last *Page) {
	p.first = first
	p.prev = prev
	p.next = next
	p.last = last
}

func (p *Page) fillMeta() {
	if p.first != nil {
		p.addMeta("start", p.first.title, p.first.Path())
	}
	if p.prev != nil {
		p.addMeta("prev", p.prev.title, p.prev.Path())
	}
	if p.next != nil {
		p.addMeta("next", p.next.title, p.next.Path())
	}
	if p.last != nil {
		p.addMeta("last", p.last.title, p.last.Path())
	}
}

func (p *Page) addMeta(rel string, title string, path string) {
	l := []string{rel, title, path}
	p.meta = append(p.meta, l)
}

func (p *Page) GetMeta() [][]string {
	p.fillMeta()
	return p.meta
}

func (p *Page) fillNavi() {
	if p.first != nil {
		p.addNavi("first", p.first.title, p.first.Path(), "&lt;&lt; first")
	}
	if p.prev != nil {
		p.addNavi("previous", p.prev.title, p.prev.Path(), "&lt; previous")
	}
	if p.next != nil {
		p.addNavi("next", p.next.title, p.next.Path(), "next &gt;")
	}
	if p.last != nil {
		p.addNavi("last", p.last.title, p.last.Path(), "newest &gt;")
	}
}

func (p *Page) IsLast() bool {
	return p.last == nil
}

func (p *Page) UrlToNext() string {
	return p.next.Path()
}

func (p *Page) addNavi(rel string, label string, title string, path string) {
	n := []string{rel, label, title, path}
	p.navi = append(p.navi, n)
}

func (p *Page) GetNavi() [][]string {
	p.fillNavi()
	return p.navi
}

func (p *Page) Path() string {
	path := config.Servedrootpath() + p.path
	return path
}

func (p *Page) FSPath() string {
	return p.path
}

func (p *Page) Img() string {
	return p.imgUrl
}
