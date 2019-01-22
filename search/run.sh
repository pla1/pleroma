#!/bin/bash
name="${PWD##*/}"
dir="$(pwd)"
echo "Build and run $name"
gofmt -w "$name".go && go build && ./search "$@"

