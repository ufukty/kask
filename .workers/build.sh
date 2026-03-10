#!/bin/bash

set -ve

git fetch --tags --quiet

if test "$WORKERS_CI_BRANCH" = "main"; then
  git checkout "$(git tag --list 'v*' --sort '-version:refname' | head -n 1)"
fi

make install
~/bin/kask build -in docs -out docs-build -domain "https://kask.ufukty.com" -v -cfw
