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
	"github.com/ingmardrewing/gomic/config"
	"github.com/ingmardrewing/gomic/page"
)

func UploadPage(p *page.Page) {

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

	stop()

	bucket := config.AwsBucket()
	localPathToFile := fmt.Sprintf("%s/%s", config.PngDir(), p.ImageFilename())
	remotePathToFile := fmt.Sprintf("%s/%s", config.AwsDir(), p.ImageFilename())
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
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", p.ImageFilename(), bucket+remotePathToFile, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", p.ImageFilename(), bucket+remotePathToFile)
	fmt.Println("Going to update db now")

	stop()
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
