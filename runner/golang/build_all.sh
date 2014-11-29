#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker build -t bazooka/runner-golang:${d%?} .
    popd
done

docker tag bazooka/runner-golang:1.3.1 bazooka/runner-golang:latest
