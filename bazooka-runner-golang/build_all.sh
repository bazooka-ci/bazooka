#!/bin/bash

for d in */ ; do
    pushd $d
    docker build -t bazooka/runner-golang:${d%?} .
    popd
done

docker tag bazooka/runner-golang:1.3.1 bazooka/runner-golang:latest
