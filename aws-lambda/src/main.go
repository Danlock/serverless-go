package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sesh = session.Must(session.NewSession())
var S3 = s3.New(sesh)
var bucketToScan = "golang-lambda-demo"

func BucketScanner() error {
	objectList, err := S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: &bucketToScan,
	})
	if err != nil {
		return err
	}

	totalCount := *objectList.KeyCount
	var totalSize int64
	for _, obj := range objectList.Contents {
		if obj.Size != nil {
			totalSize += *obj.Size
		}
	}
	log.Printf("Info about requested bucket:\nFile Count:\t%d\tTotal Size (bytes):\t%d\n", totalCount, totalSize)
	return nil
}

func main() {
	lambda.Start(BucketScanner)
}
