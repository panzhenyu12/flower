#!/bin/bash

#GO15VENDOREXPERIMENT=1 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o flower main.go
