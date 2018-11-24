#!/bin/sh

CUR_DIR=$(pwd)

if [ ! -f $CUR_DIR/certs/server.key ]; then
    cd $CUR_DIR/certs
    openssl genrsa -out server.key 2048
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
    cd $CUR_DIR
fi

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=3072 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072. -memprofile ./report/memprofile.out -cpuprofile ./report/profile.out

docker stop $CONTAINER_ID >/dev/null
