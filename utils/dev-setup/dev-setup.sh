#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "bazooka/parser" "bazooka/parsergolang" "bazooka/parserjava" \
"bazooka/orchestration" "bazooka/server" "bazooka/runnergolang" \
"bazooka/runnerjava" "bazooka/scmgit" "bazooka-cli" )

git_projects=( "bazooka" )

mkdir -p "$GOPATH/src/github.com/haklop"

for project in "${git_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop"
  if [ ! -d "$project" ]; then
    git clone git@github.com:haklop/$project.git
  else
    pushd $project
    git pull --rebase
    popd
  fi
  popd
done

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/$project"
    go get -u -v ./...
  popd
done
