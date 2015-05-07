#!/bin/bash

set -e

docker_projects=( "parser" "orchestration" "server" "web")

if [ -n "$DO_PUSH" ]; then
  docker login -e "$DOCKER_EMAIL" -p "$DOCKER_PASSWORD" -u "$DOCKER_USERNAME"

  for project in "${docker_projects[@]}"
  do
    image="bazooka/$project"
    if [ -n "$BZK_VERSION" ]; then
      docker tag "bazooka/$project" "bazooka/$project:$BZK_VERSION"
      image="bazooka/$project:$BZK_VERSION"
    fi
    docker push "$image"
  done
else
  echo "Variable DO_PUSH not defined, skipping docker push"
fi
