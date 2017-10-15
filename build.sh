#!/usr/bin/env bash

cd $(dirname $0)

./clean.sh

start=`date +%s`

go get github.com/LorisFriedel/go-bot
mkdir $(pwd)/bin

sudo docker ps > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo "Build without Docker"
    CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" github.com/LorisFriedel/go-bot
    mv $(pwd)/go-bot $(pwd)/bin/go-bot
else
    echo "Build with Docker"
    #godep save
    sudo docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.9.1 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o bin/go-bot'
fi


end=`date +%s`

runtime=$((end-start))
echo "Build done. ("${runtime}" s)"
