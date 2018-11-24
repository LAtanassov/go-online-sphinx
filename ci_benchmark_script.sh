#!/bin/sh

CUR_DIR=$(pwd)

if [ ! -f $CUR_DIR/certs/server.key ]; then
    cd $CUR_DIR/certs
    openssl genrsa -out server.key 2048
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
    cd $CUR_DIR
fi

# 32 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=32 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._32.

docker stop $CONTAINER_ID >/dev/null

# 128 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=128 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._128.

docker stop $CONTAINER_ID >/dev/null

# 512 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=512 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .512.

docker stop $CONTAINER_ID >/dev/null


# 1024 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=1024 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .1024.

docker stop $CONTAINER_ID >/dev/null


# 2048 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=2048 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .2048.

docker stop $CONTAINER_ID >/dev/null

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=3072 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072.

docker stop $CONTAINER_ID >/dev/null
