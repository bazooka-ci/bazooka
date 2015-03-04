#!/bin/bash
set -e

go_projects=( "scm/git" "scm/hg")

for project in "${go_projects[@]}"
do
  pushd "$project"
  make
  popd
done
