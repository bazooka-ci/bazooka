#!/bin/bash

set -e

docker_projects=( "parser" "orchestration" "server" "web")

for project in "${docker_projects[@]}"
do
  pushd "$project"
    make
  popd
done
