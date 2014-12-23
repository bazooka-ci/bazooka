#!/bin/bash
set -e

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

export PREFIX=$s

go_projects=( "scm/git" "scm/hg")

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
  make
  popd
done
