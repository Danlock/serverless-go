#!/bin/bash
# Installs, builds, packages and deploys to aws-lambda using AWS Cloudformation and AWS SAM

declare START_TIME=$SECONDS
function printTimeElapsed {
    echo "$(basename "$0") took $(($SECONDS - $START_TIME)) seconds to run."
}
trap printTimeElapsed EXIT

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/lambda github.com/aws/aws-lambda-go/events

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/getRandomPeople src/main.go

echo "Packaging..."
zip -j dist/getRandomPeople.zip ./dist/getRandomPeople
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation package \
--output-template-file ./dist/generated.yml \
--s3-bucket golang-lambda-demo \
--template-file template.yml

echo "Deploying to aws-lambda using cloudformation..."
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation deploy \
--capabilities CAPABILITY_NAMED_IAM \
--template-file ./dist/generated.yml \
--stack-name golang-lambda-demo-stack
