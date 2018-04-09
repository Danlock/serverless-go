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
GOOS=linux GOARCH=amd64 go build -o dist/main src/main.go

echo "Packaging..."
zip -j dist/main.zip dist/main

echo "Deploying to aws-lambda..."
AWS_DEFAULT_REGION=ca-central-1 aws \
--cli-read-timeout 0 \
--cli-connect-timeout 0 \
lambda update-function-code \
--function-name BucketScanner \
--zip-file fileb://./dist/main.zip