#!/usr/bin/env bash

set -e

if [[ ! -e "$1" ]]; then
    echo "error: folder '$1' does not exist"
    exit 1
fi

cd "$1"
npm run dev