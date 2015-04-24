#!/bin/bash

BINTRAY_PATH="https://api.bintray.com/content/bazooka/bazooka/bzk/nightly/"

versions=( "linux_amd64" "linux_386" "linux_arm" "darwin_amd64" "darwin_386")

for version in "${versions[@]}"
do
  # Upload to bintray
  curl -u "$BINTRAY_USER:$BINTRAY_API_KEY" \
    --data-binary "@bzk_$version" \
    -H "X-Bintray-Publish: 1" -H "X-Bintray-Override: 1" \
    -X PUT "$BINTRAY_PATH/bzk_$version/bzk"
done
