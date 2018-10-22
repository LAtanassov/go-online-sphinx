# go-online-sphinx

Online SPHINX - inspired by [SPHINX](https://ieeexplore.ieee.org/document/7980050)

[![Build Status](https://travis-ci.com/LAtanassov/go-online-sphinx.svg?branch=master)](https://travis-ci.com/LAtanassov/go-online-sphinx)
[![GoDoc](https://godoc.org/github.com/LAtanassov/go-online-sphinx?status.svg)](https://godoc.org/github.com/LAtanassov/go-online-sphinx)
[![Coverage Status](https://coveralls.io/repos/github/LAtanassov/go-online-sphinx/badge.svg?branch=master)](https://coveralls.io/github/LAtanassov/go-online-sphinx?branch=master)

# Rough Protocol

## Register phase

1. User registers - generate on
   - client side: k_C
   - server side: k_0(S), Q_0, k_v, \delta k_C and key material for domain passwords

## Login Phase

1. User logs in

   - client sends cID, cNonce, b (blinded password), q (group)
   - server responds sID, sNonce, bd (b with server key), Q_0 and k_v

2. Key calculation

   - client calculates mk (master key) and SKi (session key)
   - server calcualtes SKi (session key)

3. Verification
   - client sends challenge
   - server returns response

### Questions

SK_i = MAC_kv(cID | sID | cNonce | sNonce)

- cID, sID, cNonce, sNonce are public
- k_v is a secret but how is k_v shared between client and server ?
- Is MAC_kv secure => offline dictionary ?
- SK_i is session key => means session ID, session store, TLS, secure Cookies,
