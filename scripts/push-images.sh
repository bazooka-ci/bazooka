#!/bin/bash

set -e

: ${BZK_REGISTRY_EMAIL:?"BZK_REGISTRY_EMAIL has to be set."}
: ${BZK_REGISTRY_PASSWORD:?"BZK_REGISTRY_PASSWORD has to be set."}
: ${BZK_REGISTRY_USER:?"BZK_REGISTRY_USER has to be set."}
: ${BZK_REGISTRY_HOST:?"BZK_REGISTRY_HOST has to be set."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

export PREFIX=$s

$s docker login --email="$BZK_REGISTRY_EMAIL" --password="$BZK_REGISTRY_PASSWORD" --username="$BZK_REGISTRY_USER"

docker_projects=( "parser" "parserlang/golang" "parserlang/java" "orchestration" \
"server" "runner/golang" "runner/java" "scm/git" )

for project in "${docker_projects[@]}"
do
  pushd "$GOPATH/src/github.com/haklop/bazooka/$project"
  make push
  popd
done
