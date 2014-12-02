#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" "server" "cli" )

for project in "${go_projects[@]}"
do
  pushd "$project"
    godep restore ./...
  popd
done
