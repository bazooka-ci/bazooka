#!/bin/bash
set -ev

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go get -u github.com/kisielk/errcheck

go_projects=( "parser" "orchestration" "server" "cli" )

for project in "${go_projects[@]}"
do
  pushd "${GOPATH//:}/src/github.com/bazooka-ci/bazooka/$project"
    go get -t -v ./...
  popd
done
