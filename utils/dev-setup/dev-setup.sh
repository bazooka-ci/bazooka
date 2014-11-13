#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

go_projects=( "bazooka-scm-based-build/bazooka-parser" "bazooka-scm-based-build/bazooka-parser-golang" \
"bazooka-scm-based-build/bazooka-parser-java" "bazooka-scm-based-build/bazooka-orchestration" "bazooka-api/bazooka-server" \
"bazooka-api/bazooka-cli")

git_projects=( "bazooka-scm-based-build" "bazooka-api")

mkdir -p "$GOPATH/src/bitbucket.org/bywan"

for project in "${git_projects[@]}"
do
  pushd "$GOPATH/src/bitbucket.org/bywan"
  if [ ! -d "$project" ]; then
    git clone git@bitbucket.org:bywan/$project.git
  else
    pushd $project
    git pull --rebase
    popd
  fi
  popd
done

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/bitbucket.org/bywan/$project"
    go get -u -v ./...
  popd
done
