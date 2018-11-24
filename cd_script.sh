#!/bin/sh

go test ./...
go build -tags netgo ./cmd/oscli/...

go build -tags netgo ./cmd/ossvc/...

docker build -t latanassov/ossvc:0.1.0 .
docker login
docker push latanassov/ossvc:0.1.0