#!/bin/sh

# 32 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=32 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._32. -benchtime=10s

# echo $(docker logs $CONTAINER_ID)

docker stop $CONTAINER_ID >/dev/null

exit 1

# 512 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=512 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .512.

echo docker log $CONTAINER_ID

docker stop $CONTAINER_ID >/dev/null


# 1024 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=1024 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .1024.

docker stop $CONTAINER_ID >/dev/null


# 2048 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=2048 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .2048.

docker stop $CONTAINER_ID >/dev/null

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=3072 latanassov/ossrv:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072.

docker stop $CONTAINER_ID >/dev/null
