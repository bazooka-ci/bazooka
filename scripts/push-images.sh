#!/bin/bash

set -e

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "parserlang/python" "parserlang/nodejs" "orchestration" \
"server" "web")

for project in "${docker_projects[@]}"
do
  pushd "${GOPATH//:}/src/github.com/bazooka-ci/bazooka/$project"
  make push
  popd
done
