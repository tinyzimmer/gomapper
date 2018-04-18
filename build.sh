#!/bin/bash

# save start dir
startDir=$(pwd)
mkdir -p build/{bin,tmp}
mkdir -p build/usr/local/share/nmap
cp -r build/scripts build/usr/local/share/nmap/scripts

cd nmap-build
docker build -t static-binaries-nmap .
docker run \
    -v "${startDir}/build/bin":/output \
    -v "${startDir}/build/usr/local/share/nmap":/share_output \
    static-binaries-nmap

if [[ "${?}" != 0 ]]; then
    echo "Nmap compilation failed"
    exit 1
fi

cd "${startDir}"

cd build
echo -n "Compiling go binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ..
echo "Done"
docker build . -t gomapper --no-cache
cd "${startDir}"
