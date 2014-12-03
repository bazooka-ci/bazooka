#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go get -u github.com/mitchellh/gox
go get -u github.com/kisielk/errcheck

go_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" "server" "cli" )

for project in "${go_projects[@]}"
do
  pushd "../$project"
    go get -u -v ./...
  popd
done
