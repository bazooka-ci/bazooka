#!/bin/bash
# Ensure permissions are right on the key file
if [ -e "/bazooka-key" ]; then
  chmod 0600 /bazooka-key
fi

git clone "$BZK_SCM_URL" --recursive --branch "$BZK_SCM_REFERENCE" /bazooka
