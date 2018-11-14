#!/bin/sh

# 32 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=32 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._32. -benchtime=1m

docker logs $CONTAINER_ID > ./report/container_bt.log 2>&1

docker stop $CONTAINER_ID >/dev/null

# 128 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=128 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._128. -benchtime=1m

docker logs $CONTAINER_ID > ./report/container_bt.log 2>&1

docker stop $CONTAINER_ID >/dev/null

# 512 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=512 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .512. -benchtime=1m

echo docker log $CONTAINER_ID

docker stop $CONTAINER_ID >/dev/null


# 1024 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=1024 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .1024. -benchtime=1m

docker stop $CONTAINER_ID >/dev/null


# 2048 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=2048 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .2048. -benchtime=1m

docker stop $CONTAINER_ID >/dev/null

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=3072 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072. -benchtime=1m

docker stop $CONTAINER_ID >/dev/null
