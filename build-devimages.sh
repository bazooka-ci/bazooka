#!/bin/bash

set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" "server")

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    $s make devimage
  popd
done