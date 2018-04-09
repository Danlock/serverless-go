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
github.com/aws/aws-sdk-go/service/s3 \
github.com/aws/aws-sdk-go/service/sns \
github.com/aws/aws-sdk-go/aws

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/ChaosChimp src/chaos-chimp/main.go

echo "Packaging Go code..."
zip -j dist/CC.zip ./dist/ChaosChimp

echo "Packaging Cloudformation template..."
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation package \
--output-template-file ./dist/generated.yml \
--s3-bucket golang-lambda-demo \
--template-file chaos-template.yml

echo "Deploying to aws-lambda using cloudformation..."
AWS_DEFAULT_REGION=ca-central-1 aws cloudformation deploy \
--capabilities CAPABILITY_NAMED_IAM \
--template-file ./dist/generated.yml \
--stack-name chaos-chimp-lambda
