#!/bin/bash
set -e

git_log () {
  git --no-pager log --format="$1" -n 1 HEAD
}

extract_meta () {
  origin=$(git config --get remote.origin.url)
  sha1=$(git rev-parse HEAD)
  message=$(git_log %B | head -n 1)

  echo "origin: $origin"                > /meta/scm
  echo "reference: $BZK_SCM_REFERENCE"  >> /meta/scm
  echo "commit_id: $sha1"               >> /meta/scm
  echo "author:"                        >> /meta/scm
  echo "  name: $(git_log %an)"         >> /meta/scm
  echo "  email: $(git_log %ae)"        >> /meta/scm
  echo "committer:"                     >> /meta/scm
  echo "  name: $(git_log %cn)"         >> /meta/scm
  echo "  email: $(git_log %ce)"        >> /meta/scm
  echo "date: $(git_log %cd)"           >> /meta/scm
  echo "message: $message"              >> /meta/scm
}

# Ensure permissions are right on the key file
if [ -e "/bazooka-key" ]; then
  chmod 0600 /bazooka-key
fi

git clone "$BZK_SCM_URL" --recursive /bazooka
pushd /bazooka
  git checkout "$BZK_SCM_REFERENCE"
  extract_meta
popd
