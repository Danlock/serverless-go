#!/bin/bash
# Installs, builds, packages and deploys directly to aws-lambda

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/events github.com/aws/aws-lambda-go/lambda

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/main src/main.go

echo "Packaging..."
zip -j dist/main.zip dist/main

echo "Deploying to aws-lambda..."
aws lambda update-function-code \
--function-name RandomUserGenerator \
--zip-file fileb://./dist/main.zip