#!/bin/bash
# requires docker, and performs a native linux compilation of a go program (uses golang 1.8 container)
# default output file is app.out and might overidden by a program argument
name=${1:-app.out}
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.8 sh -c "CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags=\"-s\" -o ${name}"
