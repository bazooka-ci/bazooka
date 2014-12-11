#!/bin/bash
set -e

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd "$d"
      $s docker build -t "bazooka/runner-golang:${d%?}" .
    popd
done

$s docker tag bazooka/runner-golang:1.4 bazooka/runner-golang:latest
