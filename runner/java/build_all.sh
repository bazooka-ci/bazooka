#!/bin/bash

for d in */ ; do
    pushd $d
    docker build -t bazooka/runner-java:${d%?} .
    popd
done
