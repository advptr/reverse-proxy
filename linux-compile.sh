#!/bin/bash

name=$1
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.8 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o app.out'
[ "$name" != "" ] && mv app.out $1
