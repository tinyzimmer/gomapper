#!/bin/bash

set -e
set -o pipefail
set -x


NMAP_VERSION=7.70
OPENSSL_VERSION=1.0.2c


function build_openssl() {
    cd build

    # Download
    curl -LO https://www.openssl.org/source/openssl-${OPENSSL_VERSION}.tar.gz
    tar zxvf openssl-${OPENSSL_VERSION}.tar.gz
    cd openssl-${OPENSSL_VERSION}

    # Configure
    CC='/usr/bin/x86_64-alpine-linux-musl-gcc -static' ./Configure no-shared linux-x86_64

    # Build
    make
    echo "** Finished building OpenSSL"
}

function build_nmap() {
    cd /build

    # Download
    curl -LO http://nmap.org/dist/nmap-${NMAP_VERSION}.tar.bz2
    tar xjvf nmap-${NMAP_VERSION}.tar.bz2
    cd nmap-${NMAP_VERSION}

    # Configure
    CC='/usr/bin/x86_64-alpine-linux-musl-gcc -static -fPIC' \
        CXX='/usr/bin/x86_64-alpine-linux-musl-g++ -static -static-libstdc++ -fPIC' \
        LD=/usr/x86_64-alpine-linux-musl/bin/ld \
        LDFLAGS="-L/build/openssl-${OPENSSL_VERSION}"   \
        ./configure \
            --without-ndiff \
            --without-zenmap \
            --without-nmap-update \
            --with-pcap=linux \
            --with-openssl=/build/openssl-${OPENSSL_VERSION}

    # Don't build the libpcap.so file
    sed -i -e 's/shared\: /shared\: #/' libpcap/Makefile

    # Build
    make -j4
    /usr/x86_64-alpine-linux-musl/bin/strip nmap # ncat/ncat nping/nping
}

function doit() {
    uid="${1}"
    build_openssl
    build_nmap

    # Copy to output
    if [ -d /output ]
    then
        OUT_DIR=/output
        SHARE_OUT=/share_output
        mkdir -p $OUT_DIR
        cp /build/nmap-${NMAP_VERSION}/nmap $OUT_DIR/
        find /build/nmap-${NMAP_VERSION} -type f \
            -maxdepth 1 \
            -name "nmap-*" \
            -exec cp {} ${SHARE_OUT}/ \; \
            -exec /bin/echo {} \;
        cp /build/nmap-${NMAP_VERSION}/nse_main.lua "${SHARE_OUT}/" && echo "Copied nse_main"
        cp -r /build/nmap-${NMAP_VERSION}/nselib "${SHARE_OUT}/nselib" && echo "Copied nselib"
        cp -r /build/nmap-${NMAP_VERSION}/scripts "${SHARE_OUT}/scripts" && echo "Copied NSE Scripts"
        echo " ** Ensuring ownership on output **"
        chown -R ${uid}: "${OUT_DIR}" "${SHARE_OUT}"
        #cp /build/nmap-${NMAP_VERSION}/ncat/ncat $OUT_DIR/
        #cp /build/nmap-${NMAP_VERSION}/nping/nping $OUT_DIR/
        echo "** Finished **"
    else
        echo "** /output does not exist **"
    fi
}

doit "${1}"
