#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker build -t bazooka/runner-java:${d%?} .
    popd
done
