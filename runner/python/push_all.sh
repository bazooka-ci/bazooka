#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker tag bazooka/runner-python:${d%?} $BZK_REGISTRY_HOST/bazooka/runner-python:${d%?}
    $s docker push $BZK_REGISTRY_HOST/bazooka/runner-python:${d%?}
    popd
done
