FROM alpine

ARG UID
ENV UID=${UID}

RUN apk add --no-cache \
    gcc g++ musl musl-dev libc-dev \
    curl tar make autoconf m4 perl \
    libssh2 libssh2-dev linux-headers

# Add our build script
ADD build.sh /build/build.sh

# This builds the program and copies it to /output
CMD sh /build/build.sh ${UID}
