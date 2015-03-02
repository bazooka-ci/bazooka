#!/bin/bash

set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "parserlang/nodejs" "parserlang/python" "orchestration" "server")

for project in "${docker_projects[@]}"
do
  pushd "${GOPATH//:}/src/github.com/bazooka-ci/bazooka/$project"
    make
  popd
done
