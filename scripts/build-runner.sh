#!/bin/bash
set -e

go_projects=( "runner/golang" "runner/java" "runner/nodejs" "runner/python")

for project in "${go_projects[@]}"
do
  pushd "$project"
  make
  popd
done
