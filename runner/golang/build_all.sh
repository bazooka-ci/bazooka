#!/bin/bash
set -e

for d in */ ; do
    pushd "$d"
      docker build -t "bazooka/runner-golang:${d%?}" .
    popd
done

docker tag -f bazooka/runner-golang:1.4 bazooka/runner-golang:latest
