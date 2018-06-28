#!/bin/bash

source $(dirname $0)/setup-go.sh

echo "Building Releases"
make release
