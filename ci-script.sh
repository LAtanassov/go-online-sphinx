#!/bin/sh

# test and coverage
go test ./... -coverprofile=./report/coverage/coverage.out 2>&1 > ./report/test/test.txt
go tool cover -html=./report/coverage/coverage.out -o ./report/coverage/coverage.html

# allocation
go build -gcflags "-m -m" ./cmd/osctl/... 2>&1 | tee ./report/allocation/osctl.txt
go build -gcflags "-m -m" ./cmd/ossrv/... 2>&1 | tee ./report/allocation/ossrv.txt

go build -tags netgo ./cmd/osctl/...
go build -tags netgo ./cmd/ossrv/...