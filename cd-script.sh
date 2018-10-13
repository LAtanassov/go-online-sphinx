#!/bin/sh

go test ./...
go build -tags netgo ./cmd/osctl/...

go build -tags netgo ./cmd/ossrv/...

docker build -t latanassov/ossrv:0.1.0 .
docker login
docker push latanassov/ossrv:0.1.0