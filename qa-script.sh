#!/bin/sh

mkdir -p report

# race detection
go test ./... -race 2>&1 | tee ./report/race.txt

# coverage
go test ./... -coverprofile=./report/coverage.out
go tool cover -html=./report/coverage.out -o ./report/coverage.html

# allocation
go build -gcflags "-m -m" ./cmd/osctl/... 2>&1 | tee ./report/allocation_osctl.txt
go build -gcflags "-m -m" ./cmd/ossrv/... 2>&1 | tee ./report/allocation_ossrv.txt