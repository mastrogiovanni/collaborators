#!/bin/bash

docker run -v "$PWD":/root mastrogiovanni/collaborator:v0.0.1

# -e GOOS=darwin -e GOARCH=amd64 --rm  -w /app golang:1.25 go build -o generate cmd/generate/main.go

# docker run --rm -p 0.0.0.0:8080:80 -v ./:/usr/share/nginx/html nginx

