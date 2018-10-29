# go-online-sphinx

Online SPHINX - inspired by [SPHINX](https://ieeexplore.ieee.org/document/7980050)
DO NOT USE FOR PASSWORD MANAGEMENT  
THIS IS A PROTOTYPE to evaluate the Online SPHINX protocol

[![Build Status](https://travis-ci.com/LAtanassov/go-online-sphinx.svg?branch=master)](https://travis-ci.com/LAtanassov/go-online-sphinx)
[![GoDoc](https://godoc.org/github.com/LAtanassov/go-online-sphinx?status.svg)](https://godoc.org/github.com/LAtanassov/go-online-sphinx)
[![Coverage Status](https://coveralls.io/repos/github/LAtanassov/go-online-sphinx/badge.svg?branch=master)](https://coveralls.io/github/LAtanassov/go-online-sphinx?branch=master)

# Design Descision History

Goal - cloud native best practices

Therefore project was started using go-kit to embrace logging, monitoring and other DevOps layers from the start on, but those layers pollute the core functionality and introduce code that has to be maintained - better solution would be to use service mesh [istio](http://istio.io)

Another design decision made early on, because it is standard was to use REST+JSON. Marshalling and unmashalling also introduce code that has to be maintained - better solution would be to use grpc instead and generate stubs and skeletons.

# Rough Protocol

## Register phase

1. User registers by POST request, therefore it generates
   - client side: k_C
   - server side: k_0(S), Q_0, k_v, \delta k_C and key material for domain passwords

tech. note:

- rate limit needed, confirmation mail - (2 step registration process)

## Login Phase

1. User logs in

   - client sends cID, cNonce, b (blinded password), q (group)
   - server responds sID, sNonce, bd (b with server key), Q_0 and k_v

tech note:

- k_v is a secret - how is it shared with the client ?
- how does the user recognize wrong password ?

2. Key calculation

   - client calculates mk (master key) and SKi (session key)
   - server calcualtes SKi (session key)

3. Verification
   - client sends challenge
   - server returns response

tech. note:

- session via cookie over HTTPS needed - state between step 1. and step 3.
- if cookie - do we need HMAC of requests anymore ?
- Is MAC_kv secure => offline dictionary ?
