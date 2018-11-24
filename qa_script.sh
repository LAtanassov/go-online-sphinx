#!/bin/sh

mkdir -p report

go test -count=1 ./... -coverprofile=./report/coverage_ut.out -race 2>&1 | tee ./report/race_ut.txt
go tool cover -html=./report/coverage_ut.out -o ./report/coverage_ut.html

# allocation
go build -gcflags "-m -m" ./cmd/oscli/... 2>&1 | tee ./report/allocation_oscli_ut.txt
go build -gcflags "-m -m" ./cmd/ossvc/... 2>&1 | tee ./report/allocation_ossvc_ut.txt