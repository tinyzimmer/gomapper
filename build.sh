#!/bin/bash

# save start dir
startDir=$(pwd)

cd build
echo -n "Compiling go binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ..
echo "Done"
docker build .
cd "${startDir}"
