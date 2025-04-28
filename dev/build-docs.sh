#!/bin/bash

set -e

git checkout "$(git tag | grep -e '^v' | sort -r | head -n 1)"
go install ./cmd/kask

git checkout dev
kask build -in docs -out docs-build -domain / -v
