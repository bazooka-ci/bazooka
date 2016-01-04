#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "parser" "orchestration" "server" "worker")

for project in "${go_projects[@]}"
do
  pushd "${GOPATH//:}/src/github.com/bazooka-ci/bazooka/$project"
    errcheck -ignore 'Close|[wW]rite.*|Encode|Flush|Seek|[rR]ead.*' github.com/bazooka-ci/bazooka/$project/...
  popd
done
