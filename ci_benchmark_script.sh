#!/bin/sh

# 32 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSVC_KEYLENGTH=32 latanassov/:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._32.

docker stop $CONTAINER_ID >/dev/null

# 128 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSVC_KEYLENGTH=128 latanassov/:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench ._128.

docker stop $CONTAINER_ID >/dev/null

# 512 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSVC_KEYLENGTH=512 latanassov/:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .512.

docker stop $CONTAINER_ID >/dev/null


# 1024 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSVC_KEYLENGTH=1024 -e OSSVC_TIMEOUTSEC=60 latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .1024.

docker stop $CONTAINER_ID >/dev/null


# 2048 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e ossvc_KEYLENGTH=2048 -e OSSVC_TIMEOUTSEC=60 latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .2048.

docker stop $CONTAINER_ID >/dev/null

# 3072 bit key length
CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSVC_KEYLENGTH=3072 -e OSSVC_TIMEOUTSEC=60 latanassov/ossvc:0.1.0)

go test -benchmem -run None github.com/LAtanassov/go-online-sphinx/pkg/client -bench .3072.

docker stop $CONTAINER_ID >/dev/null
