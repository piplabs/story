#!/usr/bin/env bash

# Runs go test for all touched packages
MOD=$(go list -m)
PKGS=$(echo "$@"| xargs -n1 dirname | sort -u | sed -e "s#^#${MOD}/#")

# TODO: fix tests and enable
# go test -tags=verify_logs -failfast -race -timeout=2m $PKGS
