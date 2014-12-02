#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" "server")

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    errcheck -ignore 'Close|[wW]rite.*|Encode|Flush|Seek|[rR]ead.*' github.com/haklop/bazooka/$project/...
  popd
done
