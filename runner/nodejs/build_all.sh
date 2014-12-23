#!/bin/bash
set -e

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd "$d"
      $s docker build -t "bazooka/runner-nodejs:${d%?}" .
    popd
done
