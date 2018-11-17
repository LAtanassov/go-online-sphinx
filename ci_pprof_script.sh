#!/bin/sh

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=3072 -e OSSRV_TIMEOUTSEC=60 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072. -memprofile ./report/memprofile.out -cpuprofile ./report/profile.out

docker stop $CONTAINER_ID >/dev/null
