#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

go_projects=( "bazooka/parser" "bazooka/parsergolang" "bazooka/parserjava" \
"bazooka/orchestration" "bazooka/server" "bazooka/runnergolang" \
"bazooka/runnerjava" "bazooka/scmgit")

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/$project"
    $s make
  popd
done
