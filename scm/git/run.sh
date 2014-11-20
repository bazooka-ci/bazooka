#!/bin/bash
# Ensure permissions are right on the key file
set -e

if [ -e "/bazooka-key" ]; then
  chmod 0600 /bazooka-key
fi

git clone "$BZK_SCM_URL" --recursive /bazooka
pushd /bazooka
  git checkout "$BZK_SCM_REFERENCE"
  git rev-parse HEAD > /meta/reference
popd
