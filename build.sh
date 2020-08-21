#!/bin/bash

os="$(uname -s)"

test -d builds || mkdir builds
for i in linux darwin; do
env GOOS=$i GOARCH=amd64 go build -o builds/vpngater-$i .
done

