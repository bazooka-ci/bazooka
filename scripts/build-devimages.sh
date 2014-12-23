#!/bin/bash

set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

export PREFIX=$s

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "parserlang/python" "parserlang/nodejs" "orchestration" "server")

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    make devimage
  popd
done
