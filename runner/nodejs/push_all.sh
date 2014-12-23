#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker tag bazooka/runner-nodejs:${d%?} $BZK_REGISTRY_HOST/bazooka/runner-nodejs:${d%?}
    $s docker push $BZK_REGISTRY_HOST/bazooka/runner-nodejs:${d%?}
    popd
done
