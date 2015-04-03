#!/bin/bash

set -e

docker_projects=( "parser" "orchestration" "server" "web")

docker login -e "$DOCKER_EMAIL" -p "$DOCKER_PASSWORD" -u "$DOCKER_USERNAME"

for project in "${docker_projects[@]}"
do
  docker push "bazooka/$project"
done
