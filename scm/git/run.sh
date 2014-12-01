#!/bin/bash
set -e

function extract_meta {
  sha1=$(git rev-parse HEAD)
  author_name=$(git --no-pager log --format='%an' -n 1 HEAD)
  author_email=$(git --no-pager log --format='%ae' -n 1 HEAD)
  date=$(git --no-pager log --format='%cd' -n 1 HEAD)
  message=$(git --no-pager log --format=%B -n 1 | head -n 1)

  echo "reference: $BZK_SCM_REFERENCE"  > /meta/scm
  echo "commit_id: $sha1"               >> /meta/scm
  echo "author:"                        >> /meta/scm
  echo "  name: $author_name"           >> /meta/scm
  echo "  email: $author_email"         >> /meta/scm
  echo "date: $date"                    >> /meta/scm
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
