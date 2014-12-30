#!/bin/bash
set -e

for d in */ ; do
    pushd "$d"
      docker build -t "bazooka/runner-java:${d%?}" .
    popd
done

docker tag -f bazooka/runner-java:oraclejdk8 bazooka/runner-java:latest
