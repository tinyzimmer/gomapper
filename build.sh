#!/bin/bash

# save start dir
startDir=$(pwd)
mkdir -p build/{bin,tmp}
mkdir -p build/usr/local/share/nmap

echo "Starting static nmap build"
cd nmap-build
docker build -t nmap-build --build-arg UID=$(id -u) .
docker run --rm \
    -e UID=$(id -u) \
    -v "${startDir}/build/bin":/output \
    -v "${startDir}/build/usr/local/share/nmap":/share_output \
    nmap-build

if [[ "${?}" != 0 ]]; then
    echo "Nmap compilation failed"
    exit 1
fi

cd "${startDir}"

cd build
echo -n "Compiling static go binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ..
echo "Done"
upx bin/main
docker build . -t gomapper --no-cache
cd "${startDir}"
