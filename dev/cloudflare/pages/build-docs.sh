#!/bin/bash

set -ve

git fetch --tags --quiet
git checkout "$(git tag --list 'v*' | sort -Vr | head -n 1)"
make install

git checkout main --
~/bin/kask build -in docs -out docs-build -domain / -v
