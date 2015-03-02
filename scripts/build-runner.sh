#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "runner/golang" "runner/java" "runner/nodejs" "runner/python")

for project in "${go_projects[@]}"
do
  pushd "${GOPATH//:}/src/github.com/bazooka-ci/bazooka/$project"
  make
  popd
done
