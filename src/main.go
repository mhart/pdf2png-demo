package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	os.Chdir("/tmp") // Change working directory to /tmp so we can write files

	if len(os.Args) > 1 && os.Args[1] == "--test" {
		pdf2png()
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, s3Event events.S3Event) ([]string, error) {
	sess := session.Must(session.NewSession())

	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key

	if !strings.HasSuffix(strings.ToLower(key), ".pdf") {
		return nil, fmt.Errorf("Key does not look like a PDF file: %s", key)
	}
	convertKey := "converted/" + strings.TrimPrefix(key[:len(key)-4], "upload/") // trim ".pdf" off the end

	fmt.Println("Received S3 event, Bucket: ", bucket, ", Key: ", key, "Converting to: ", convertKey)

	// Will download the S3 object locally to /tmp/input.pdf
	if err := downloadPdf(sess, bucket, key); err != nil {
		return nil, err
	}

	// Convert to PNGs
	if err := pdf2png(); err != nil {
		return nil, err
	}

	// Upload all PNGs from /tmp/output to our converted/ key prefix
	return uploadPngs(sess, bucket, convertKey)
}

func pdf2png() error {
	cmds := `
	rm -rf output
	mkdir output
	pdftocairo -png input.pdf output/converted
	`

	cmd := exec.Command("sh", "-cex", strings.TrimSpace(cmds))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadPdf(sess *session.Session, bucket string, key string) error {
	inputFile, err := os.Create("input.pdf")
	if err != nil {
		return err
	}
	defer inputFile.Close()

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(inputFile, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	return err
}

func uploadPngs(sess *session.Session, bucket string, keyPrefix string) ([]string, error) {
	var wg sync.WaitGroup

	files, err := ioutil.ReadDir("output")
	if err != nil {
		return nil, err
	}
	wg.Add(len(files))
	s3Locations := make([]string, len(files))

	uploader := s3manager.NewUploader(sess)

	var uploadErr error

	for ix, file := range files {
		go func(filename string, ix int) {
			defer wg.Done()

			pngfile, _ := os.Open(filepath.Join("output", filename))
			defer pngfile.Close()

			result, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: &bucket,
				Key:    aws.String(keyPrefix + "/" + filename),
				Body:   pngfile,
			})
			if err != nil {
				uploadErr = err
			} else {
				s3Locations[ix] = result.Location
			}
		}(file.Name(), ix)
	}

	wg.Wait()

	return s3Locations, uploadErr
}
