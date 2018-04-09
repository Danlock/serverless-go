#!/bin/bash
# Installs, builds, packages and deploys directly to aws-lambda

declare START_TIME=$SECONDS
function printTimeElapsed {
    echo "$(basename "$0") took $(($SECONDS - $START_TIME)) seconds to run."
}
trap printTimeElapsed EXIT

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/lambda \
    github.com/aws/aws-sdk-go/service/s3 \
    github.com/aws/aws-sdk-go/aws/session

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/BucketScanner src/main.go

echo "Packaging..."
zip -j dist/BS.zip ./dist/BucketScanner

echo "Deploying to aws-lambda..."
AWS_DEFAULT_REGION=ca-central-1 aws \
--cli-read-timeout 0 \
--cli-connect-timeout 0 \
lambda create-function \
--function-name BucketScanner \
--runtime "go1.x" \
--role "arn:aws:iam::736759784672:role/service-role/testLambdaRole" \
--handler BucketScanner \
--zip-file fileb://./dist/BS.zip