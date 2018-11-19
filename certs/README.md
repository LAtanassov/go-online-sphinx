# generate test key and cert

Key considerations for algorithm "RSA" ≥ 2048-bit

```sh
openssl genrsa -out server.key 2048
```

Key considerations for algorithm "ECDSA" ≥ secp384r1
List ECDSA the supported curves (openssl ecparam -list_curves)

```sh
openssl ecparam -genkey -name secp384r1 -out server.key
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```sh
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
