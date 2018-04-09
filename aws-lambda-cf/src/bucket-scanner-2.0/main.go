package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sesh = session.Must(session.NewSession())
var S3 = s3.New(sesh)
var bucketToScan = os.Getenv("BUCKET_TO_SCAN")

func BucketScannerV2(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	listObjInput := s3.ListObjectsV2Input{Bucket: &bucketToScan}
	keyToScan := req.QueryStringParameters["key"]
	if len(keyToScan) > 0 {
		listObjInput.Prefix = &keyToScan
	}

	objectList, err := S3.ListObjectsV2(&listObjInput)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	totalCount := *objectList.KeyCount
	if totalCount < 1 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	var totalSize int64
	for _, obj := range objectList.Contents {
		if obj.Size != nil {
			totalSize += *obj.Size
		}
	}

	log.Printf("Info about requested bucket:\nFile Count:\t%d\tTotal Size (bytes):\t%d\n", totalCount, totalSize)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"fileCount": %d, "sizeInBytes": %d}`, totalCount, totalSize),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(BucketScannerV2)
}
