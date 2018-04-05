#!/bin/bash
# Installs, builds, packages and deploys directly to aws-lambda

declare START_TIME=$SECONDS
function printTimeElapsed {
    echo "$(basename "$0") took $(($SECONDS - $START_TIME)) seconds to run."
}
trap printTimeElapsed EXIT

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/lambda

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/main src/main.go

echo "Packaging..."
zip -j dist/main.zip ./dist/main

echo "Deploying to aws-lambda..."
AWS_DEFAULT_REGION=ca-central-1 aws lambda create-function \
--function-name RandomUserGenerator \
--runtime "go1.x" \
--role "arn:aws:iam::736759784672:role/service-role/testLambdaRole" \
--handler main \
--zip-file fileb://./dist/main.zip