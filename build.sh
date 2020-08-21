#!/bin/bash

test -d builds || mkdir builds
for i in linux darwin; do
	docker run --rm -d -v $PWD:/go/src/vpngater -w /go/src/vpngater golang:1.14-buster env GOOS=$i GOARCH=amd64 go build -o builds/vpngater-$i
done