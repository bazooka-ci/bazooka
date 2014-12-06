#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

export PREFIX=$s

go_projects=( "runner/golang" "runner/java" )

for project in "${go_projects[@]}"
do
  pushd "../$project"
  make
  popd
done
