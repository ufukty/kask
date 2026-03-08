#!/bin/bash

set -ve

printenv | sort

git fetch --tags --quiet
if test "$WORKERS_CI_BRANCH" = "main"; then
  git checkout "$(git tag --list 'v*' | sort -Vr | head -n 1)"
  make install
  git checkout main --
  ~/bin/kask build -in docs -out docs-build -domain "https://kask.ufukty.com" -v -cfw
else
  make install
  ~/bin/kask build -in docs -out docs-build -domain "/" -v -cfw
fi
