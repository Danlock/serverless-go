#!/bin/bash
# Installs, builds, packages and deploys to aws-lambda using AWS Cloudformation and AWS SAM

declare START_TIME=$SECONDS
function printTimeElapsed {
    echo "$(basename "$0") took $(($SECONDS - $START_TIME)) seconds to run."
}
trap printTimeElapsed EXIT

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/lambda \
github.com/aws/aws-lambda-go/events \
github.com/aws/aws-sdk-go/aws/session \
github.com/aws/aws-sdk-go/service/s3

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/BucketScannerV2 src/bucket-scanner-2.0/main.go

echo "Packaging Go code..."
zip -j dist/BSV2.zip ./dist/BucketScannerV2

echo "Packaging Cloudformation template..."
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation package \
--output-template-file ./dist/generated.yml \
--s3-bucket golang-lambda-demo \
--template-file bucket-template.yml

echo "Deploying to aws-lambda using cloudformation..."
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation deploy \
--capabilities CAPABILITY_NAMED_IAM \
--template-file ./dist/generated.yml \
--stack-name bucket-scanner-lambda
