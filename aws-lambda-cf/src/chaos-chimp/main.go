package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
)

var sesh = session.Must(session.NewSession())
var S3 = s3.New(sesh)

var snsRegion = "us-east-1"
var SNS = sns.New(sesh, &aws.Config{Region: &snsRegion})

var doomedBucket = os.Getenv("DOOMED_BUCKET")
var doomedFolder = os.Getenv("DOOMED_FOLDER")

func NotifyVictim(phoneNum, deletedFile string) error {
	msg := fmt.Sprintf("Uh oh! Something happened to %s... - Chaos Chimp \U0001f412", deletedFile)
	_, err := SNS.Publish(&sns.PublishInput{Message: &msg, PhoneNumber: &phoneNum})
	return err
}

func ChaosChimp(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	phoneNumber := req.QueryStringParameters["cell"]
	if len(phoneNumber) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"msg":"Must provide a E164 phone number!"}`,
		}, nil
	}

	objectList, err := S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: &doomedBucket,
		Prefix: &doomedFolder,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	fileAmt := len(objectList.Contents)
	if fileAmt == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusExpectationFailed,
			Body:       `{"msg":"Doomed folder is empty! Someone's slacking..."}`,
		}, nil
	}

	luckyFileIdx := rand.Intn(fileAmt)
	luckyFile := objectList.Contents[luckyFileIdx]

	if err := NotifyVictim(phoneNumber, *luckyFile.Key); err != nil {
		log.Printf("Err: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusServiceUnavailable,
			Body:       `{"msg":"SNS failed to send! Must provide an E164 phone number..."}`,
		}, nil
	}

	_, err = S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &doomedBucket,
		Key:    luckyFile.Key,
	})
	if err != nil {
		log.Printf("Err: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"msg":"Failed to delete anything! The victim can breathe easy..."}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       `{"msg": "Chaos has been created."}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(ChaosChimp)
}
