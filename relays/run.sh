#!/bin/bash
name="${PWD##*/}"
echo "Build and run $name"
gofmt -w "$name".go && go build && ./"$name"

