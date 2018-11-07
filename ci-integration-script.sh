#!/bin/sh

CONTAINER_ID=$(docker run -d -p 8080:8080 latanassov/ossrv:0.1.0)

go test -count=1 -tags=integration ./...

docker stop $CONTAINER_ID >/dev/null