# go-online-sphinx

Online SPHINX - inspired by [SPHINX](https://ieeexplore.ieee.org/document/7980050)   
DO NOT USE FOR PASSWORD MANAGEMENT  
THIS IS A PROTOTYPE to evaluate the Online SPHINX protocol

[![Build Status](https://travis-ci.com/LAtanassov/go-online-sphinx.svg?branch=master)](https://travis-ci.com/LAtanassov/go-online-sphinx)
[![GoDoc](https://godoc.org/github.com/LAtanassov/go-online-sphinx?status.svg)](https://godoc.org/github.com/LAtanassov/go-online-sphinx)
[![Coverage Status](https://coveralls.io/repos/github/LAtanassov/go-online-sphinx/badge.svg?branch=master)](https://coveralls.io/github/LAtanassov/go-online-sphinx?branch=master)

# CI/CD/QA
```sh
>./ci_script.sh               # run unit tests & build 
>./ci_integration_script.sh   # run integration tests, assumes docker
>./ci_benchmark_script.sh     # run benchmark tests, assumes docker
>./ci_pprof_script.sh         # run benchmark tests with pprof, assumes docker

>./cd_script.sh               # run unit tests & build & docker build and push -> openshift

>./.qa_script.sh              # generates code coverage report, memory allocation
```

# Design Descision History

Goal - cloud native best practices

Therefore project was started using go-kit to embrace logging, monitoring and other DevOps layers from the start on, but those layers pollute the core functionality and introduce code that has to be maintained - better solution would be to use service mesh [istio](http://istio.io)

Another design decision made early on, because it is standard was to use REST+JSON. Marshalling and unmashalling also introduce code that has to be maintained - better solution would be to use grpc instead and generate stubs and skeletons.

# Online SPHINX Protocol

will be explained in details at some point.

## tech. note:

- session via cookie over HTTPS needed - stateful protocol.
- if cookie - do we need HMAC of requests anymore ?
- Is MAC_kv secure => offline dictionary ?
