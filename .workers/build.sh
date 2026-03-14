#!/bin/bash

set -ve

git fetch --tags --quiet

if test "$WORKERS_CI_BRANCH" = "main"; then
  git checkout "$(git tag --list 'v*' --sort '-version:refname' | head -n 1)"
fi

go install ./cmd/kask
go env
echo "$PATH" | sed 's/:/\n/g' | sort
PATH="$(go env GOBIN):$(go env GOPATH)/bin:${PATH}"
kask build -in docs -out docs-build -domain "https://kask.ufukty.com" -v -cfw
