#!/bin/bash
mkdir data

openssl genrsa -out data/server.key 2048
openssl ecparam -genkey -name secp384r1 -out data/server.key
openssl req -new -x509 -sha256 -key data/server.key -out data/server.crt -days 3650