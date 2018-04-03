#!/bin/bash
# Installs, builds, packages and deploys directly to aws-lambda

echo "Go getting dependencies..."
go get github.com/aws/aws-lambda-go/events github.com/aws/aws-lambda-go/lambda

echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/main src/main.go

echo "Packaging..."
zip -j dist/main.zip ./dist/main

echo "Deploying to aws-lambda..."
aws lambda create-function \
--function-name RandomUserGenerator \
--runtime "go1.x" \
--role "arn:aws:iam::736759784672:role/Sphinx-WebAppFirewall-LambdaRoleLogParser-9D643AQEH1X0" \
--handler main \
--zip-file fileb://./dist/main.zip