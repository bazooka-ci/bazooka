#!/bin/bash
set -e

hg_log () {
  hg log  -T $1 -l 1 -r .
}

extract_meta () {
  desc=$(hg log -l 1 -r . -T '{splitlines(desc)%"  {line}\n"}')

  echo "origin: $BZK_SCM_URL"                  >  /meta/scm
  echo "reference: $BZK_SCM_REFERENCE"         >> /meta/scm
  echo "commit_id: $(hg_log '{node}')"         >> /meta/scm
  echo "author:"                               >> /meta/scm
  echo "  name: $(hg_log '{author|person}')"   >> /meta/scm
  echo "  email: $(hg_log '{author|email}')"   >> /meta/scm
  echo "committer:"                            >> /meta/scm
  echo "  name: $(hg_log '{author|person}')"   >> /meta/scm
  echo "  email: $(hg_log '{author|email}')"   >> /meta/scm
  echo "date: $(hg_log '{date|rfc822date}')"   >> /meta/scm
  echo "message: |"                            >> /meta/scm
  echo "$desc"                                 >> /meta/scm
}

# Ensure permissions are right on the key file
if [ -e "/bazooka-key" ]; then
  chmod 0600 /bazooka-key
fi

hg clone "$BZK_SCM_URL" /bazooka
pushd /bazooka
  hg update "$BZK_SCM_REFERENCE"
  extract_meta
popd
