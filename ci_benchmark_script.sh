#!/bin/sh

CONTAINER_ID=$(docker run -d -p 8080:8080 latanassov/ossrv:0.1.0)

go test -benchmem -run=^$ github.com/LAtanassov/go-online-sphinx/pkg/client -bench .

docker stop $CONTAINER_ID >/dev/null

