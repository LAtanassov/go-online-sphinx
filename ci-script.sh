#!/bin/sh

go test ./...
go build -tags netgo ./cmd/osctl/...
go build -tags netgo ./cmd/ossrv/...