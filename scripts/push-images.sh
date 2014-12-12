#!/bin/bash

set -e

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "parserlang/nodejs" "orchestration" \
"server" "runner/golang" "runner/java" "scm/git" "web")

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
  make push
  popd
done
