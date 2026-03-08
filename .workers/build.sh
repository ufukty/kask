#!/bin/bash

set -ve

printenv | sort

git fetch --tags --quiet
if [[ "$WORKERS_CI_BRANCH" = "main" && "$CI_ENVIRONMENT" = "prod" ]]; then
  git checkout "$(git tag --list 'v*' | sort -Vr | head -n 1)"
  make install
  ~/bin/kask build -in docs -out docs-build -domain "https://kask.ufukty.com" -v -cfw
else
  make install
  ~/bin/kask build -in docs -out docs-build -domain "/" -v -cfw
fi
