#!/bin/bash

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

for d in */ ; do
    pushd $d
    $s docker tag bazooka/runner-java:${d%?} $(BZK_REGISTRY_HOST)/bazooka/runner-java:${d%?}
    $s docker push $(BZK_REGISTRY_HOST)/bazooka/runner-java:${d%?}
    popd
done

docker push bazooka/runner-java:latest $(BZK_REGISTRY_HOST)/bazooka/runner-java:latest
