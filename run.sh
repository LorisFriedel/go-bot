#!/usr/bin/env bash

cd $(dirname $0)
BIN_PATH=./bin

if [ ! -f $BIN_PATH/go-bot ]; then
    ./build.sh
fi

$BIN_PATH/go-bot "$@"
