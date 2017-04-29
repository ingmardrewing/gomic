package aws

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
)

func UploadPage(p comic.AwsPage) {
	sess := getAwsSession()
	uploadComicPageFile(p, sess)
	uploadComicPageThumbnailFile(p, sess)
	fmt.Println("Going to update db now")
	stop()
}

func uploadComicPageThumbnailFile(p comic.AwsPage, sess *session.Session) {
	bucket := config.AwsBucket()
	localPathToThumbnail, remotePathToThumbnail := getThumbnailPaths(p)
	upload(localPathToThumbnail, remotePathToThumbnail, sess, bucket)
}

func getThumbnailPaths(p comic.AwsPage) (string, string) {
	localPathToThumbnail := fmt.Sprintf("%sthumb_%s", config.PngDir(), p.ImageFilename())
	remotePathToThumbnail := fmt.Sprintf("%s/thumb_%s", config.AwsDir(), p.ImageFilename())
	return localPathToThumbnail, remotePathToThumbnail
}

func uploadComicPageFile(p comic.AwsPage, sess *session.Session) {
	bucket := config.AwsBucket()
	localPathToFile := fmt.Sprintf("%s%s", config.PngDir(), p.ImageFilename())
	remotePathToFile := fmt.Sprintf("%s/%s", config.AwsDir(), p.ImageFilename())
	upload(localPathToFile, remotePathToFile, sess, bucket)
}

func getAwsSession() *session.Session {
	// Initialize a session reading env vars (reading done by AWS)
	sess, _ := session.NewSession()
	return sess
}

func getS3(sess *session.Session) *s3.S3 {
	// Create S3 service client
	svc := s3.New(sess)
	return svc
}

func upload(from string, to string, sess *session.Session, bucket string) {
	file, err := os.Open(from)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(to),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", from, bucket+to, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", from, bucket+to)
}

func stop() {
	answer := AskUser("Proceed? [yN]")
	if answer {
		fmt.Println("continuing")
	} else {
		os.Exit(0)
	}
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
