package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ingmardrewing/gomic/config"
)

type AwsPage interface {
	ImageFilename() string
}

func UploadPage(p AwsPage) {
	sess := getAwsSession()

	tl, tr := getThumbnailPaths(p)
	upload(config.AwsBucket(), tl, tr, sess)

	fl, fr := getFilePaths(p)
	upload(config.AwsBucket(), fl, fr, sess)
}

func getThumbnailPaths(p AwsPage) (string, string) {
	localPathToThumbnail := fmt.Sprintf("%sthumb_%s", config.PngDir(), p.ImageFilename())
	remotePathToThumbnail := fmt.Sprintf("%s/thumb_%s", config.AwsDir(), p.ImageFilename())
	return localPathToThumbnail, remotePathToThumbnail
}

func getFilePaths(p AwsPage) (string, string) {
	localPathToFile := fmt.Sprintf("%s%s", config.PngDir(), p.ImageFilename())
	remotePathToFile := fmt.Sprintf("%s/%s", config.AwsDir(), p.ImageFilename())
	return localPathToFile, remotePathToFile
}

func getAwsSession() *session.Session {
	// Initialize a session reading env vars (reading done by AWS)
	sess, _ := session.NewSession()
	return sess
}

func upload(bucket string, from string, to string, sess *session.Session) {
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
