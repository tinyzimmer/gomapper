#!/bin/bash

# see if we can shed a few lbs
UPX_ENABLED=1
which upx &> /dev/null || UPX_ENABLED=0

startDir=$(pwd)

function main() {
    setup_build_directories
    build_third_party
    build_gomapper
}

function setup_build_directories() {
    echo "Initializing build directories"
    rm -rf build/{bin,tmp,usr}
    mkdir -p build/{bin,tmp}
    mkdir -p build/usr/local/share/nmap
}

function build_third_party() {
    echo "Starting static nmap build"
    cd build/nmap-build
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
}

function build_gomapper() {
    cd build
    echo -n "Compiling static go binary..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ..
    echo "Done"
    if [[ ${UPX_ENABLED} == 1 ]] ; then
        upx bin/main
    fi
    docker build . -t gomapper --no-cache
    cd "${startDir}"
}

if [[ "$(basename ${0})" == "build.sh" ]] ; then main ; fi
