#!/bin/bash

rm -r build

mkdir -p build/dist

GOOS=linux GOARCH=amd64 go build -o build/bootstrap cmd/forgotpassword-func/main.go

zip build/dist/app.zip build/bootstrap
