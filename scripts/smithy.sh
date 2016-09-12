#!/usr/bin/env bash

# This is so `tee` doesn't absorb a non-zero exit code
set -eo pipefail

echo $PWD

mkdir -p $SMITHY_ROOT/test_results/

# Move godeps to gopath
cp -r $FRUGAL_HOME/Godeps/_workspace/* $GOPATH/
cp -r $FRUGAL_HOME/lib/go/Godeps/_workspace/* $GOPATH/

# Run each language build and tests in parallel
cd $FRUGAL_HOME
go run scripts/smithy/parallel_smithy.go

