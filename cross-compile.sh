#!/bin/bash

env GOOS=linux GOARCH=386 go build -o ./binaries/go-couchdb-compactor-linux-386 compactor.go
env GOOS=linux GOARCH=amd64 go build -o ./binaries/go-couchdb-compactor-linux-amd64 compactor.go
env GOOS=darwin GOARCH=386 go build -o ./binaries/go-couchdb-compactor-mac-386 compactor.go
