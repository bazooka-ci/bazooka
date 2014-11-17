#!/bin/bash
set -e

go_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" "server" "cli" )

for project in "${go_projects[@]}"
do
  pushd "$project"
    go get -u -v ./...
  popd
done
