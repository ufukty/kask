#!/bin/bash

set -ve

git fetch --tags --quiet
make install
~/bin/kask build -in docs -out docs-build -domain / -cfw
