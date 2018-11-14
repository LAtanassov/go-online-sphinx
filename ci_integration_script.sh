#!/bin/sh

CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=8 latanassov/ossrv:0.1.0)

mkdir -p report

go test -count=1 -tags=integration ./... -coverprofile=./report/coverage_it.out -race 2>&1 | tee ./report/out_it.txt
go tool cover -html=./report/coverage_it.out -o ./report/coverage_it.html

docker logs $CONTAINER_ID > ./report/container_it.log 2>&1

docker stop $CONTAINER_ID >/dev/null