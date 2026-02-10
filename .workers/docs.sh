#!/bin/bash

set -ve

make install
~/bin/kask build -in docs -out docs-build -domain / -v
