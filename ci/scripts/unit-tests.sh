#!/bin/bash

source $(dirname $0)/setup-go.sh

go test -cover ./...
