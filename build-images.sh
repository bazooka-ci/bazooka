#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" \
"server" "runner/golang" "runner/java" "scm/git" )

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
    $s make
  popd
done
