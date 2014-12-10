#!/bin/bash
set -e

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd "$d"
      $s docker build -t "bazooka/runner-java:${d%?}" .
    popd
done

docker tag bazooka/runner-java:oraclejdk8 bazooka/runner-java:latest
