#!/bin/bash

openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/CN=MyLocalCA"

openssl genrsa -out payment.key 2048
openssl req -new -key payment.key -out payment.csr -subj "/CN=payment-service"

cat > payment.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = payment-service
IP.1 = 127.0.0.1
EOF

openssl x509 -req -in payment.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out payment.crt -days 365 -sha256 -extfile payment.ext

openssl genrsa -out general-client.key 2048
openssl req -new -key general-client.key -out general-client.csr -subj "/CN=internal-microservice"
openssl x509 -req -in general-client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out general-client.crt -days 365 -sha256

rm *.csr *.ext
