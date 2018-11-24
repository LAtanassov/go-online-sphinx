#!/bin/sh

CUR_DIR=$(pwd)

if [ ! -f $CUR_DIR/certs/server.key ]; then
    cd $CUR_DIR/certs
    openssl genrsa -out server.key 2048
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
    cd $CUR_DIR
fi

CONTAINER_ID=$(docker run -d -p 443:443 -e OSSVC_KEYLENGTH=8 -v $CUR_DIR/certs:/certs latanassov/ossvc:0.1.0)

mkdir -p report

go test -count=1 -tags=integration ./... -coverprofile=./report/coverage_it.out -race 2>&1 | tee ./report/out_it.txt
go tool cover -html=./report/coverage_it.out -o ./report/coverage_it.html

docker logs $CONTAINER_ID > ./report/container_it.log 2>&1

docker stop $CONTAINER_ID >/dev/null