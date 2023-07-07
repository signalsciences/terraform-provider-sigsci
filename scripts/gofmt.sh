#!/bin/bash

need_gofmt=$(gofmt -s -l .)

if [[ -n ${need_gofmt} ]]; then
    echo "The following files fail gofmt -s:"
    echo "${need_gofmt}"
    exit 1
fi
