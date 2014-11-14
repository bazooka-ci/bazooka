#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "parser" "parser/golang" "parser/java" "orchestration" "server" "cli" )

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    go get -u -v ./...
  popd
done
