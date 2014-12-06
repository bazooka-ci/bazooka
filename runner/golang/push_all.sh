#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker tag bazooka/runner-golang:${d%?} $BZK_REGISTRY_HOST/bazooka/runner-golang:${d%?}
    $s docker push $BZK_REGISTRY_HOST/bazooka/runner-golang:${d%?}
    popd
done

$s docker push bazooka/runner-golang:latest $BZK_REGISTRY_HOST/bazooka/runner-golang:latest
