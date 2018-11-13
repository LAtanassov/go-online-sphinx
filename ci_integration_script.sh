#!/bin/sh

CONTAINER_ID=$(docker run -d -p 8080:8080 -e OSSRV_KEYLENGTH=8 latanassov/ossrv:0.1.0)

go test -count=1 -tags=integration ./...

docker logs $CONTAINER_ID > container.log 2>&1

docker stop $CONTAINER_ID >/dev/null