#!/bin/bash

: ${GOPATH:?"GOPATH has to be set. See https://golang.org/doc/code.html#GOPATH for more information."}

if [ "$(uname)" != "Darwin" ]; then
  s=sudo
fi

go_projects=( "bazooka-scm-based-build/bazooka-parser" "bazooka-scm-based-build/bazooka-parser-golang" \
"bazooka-scm-based-build/bazooka-parser-java" "bazooka-scm-based-build/bazooka-orchestration" "bazooka-api/bazooka-server" \
"bazooka-scm-based-build/bazooka-runner-golang" "bazooka-scm-based-build/bazooka-runner-java" \
"bazooka-scm-based-build/bazooka-scm-git")

for project in "${go_projects[@]}"
do
  pushd "$GOPATH/src/bitbucket.org/bywan/$project"
    $s make
  popd
done
