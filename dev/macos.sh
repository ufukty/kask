#!/usr/local/bin/bash

set -ex

which dot ||
  brew install graphviz

which godepgraph ||
  go install github.com/kisielk/godepgraph@latest
