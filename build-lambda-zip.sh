#!/bin/bash
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap *.go
zip -r bootstrap.zip bootstrap html static