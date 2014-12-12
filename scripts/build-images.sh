#!/bin/bash

set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" \
"server" "runner/golang" "runner/java" "scm/git" )

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    make
  popd
done
