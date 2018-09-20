#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main build/redis-transmission
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
mv main.exe build/redis-transmission.exe
