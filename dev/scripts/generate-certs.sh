#!/bin/bash

: "${1?'missing cert directory'}"

CERTS_DIR="$1"

mkdir -p "$CERTS_DIR"
chmod 0700 "$CERTS_DIR"
#cd "$CERTS_DIR"

cat > "$CERTS_DIR"/csr.conf <<- END
[req]
default_bits = 2048
distinguished_name = dn
req_extensions     = req_ext
prompt             = no

[dn]
CN="host.docker.internal"

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = host.docker.internal
END

# ca cert and private key
openssl req -nodes -new -x509 -keyout "$CERTS_DIR"/ca.key -out "$CERTS_DIR"/ca.pem -subj "/CN=Admission Controller Demo Certificate Authority"
# server key
openssl genrsa -out "$CERTS_DIR"/server.key 2048
# server csr + cert
openssl req -new -key "$CERTS_DIR"/server.key -config "$CERTS_DIR"/csr.conf \
    | openssl x509 -req -CA "$CERTS_DIR"/ca.pem -CAkey "$CERTS_DIR"/ca.key -CAcreateserial -extfile "$CERTS_DIR"/csr.conf -extensions req_ext -out "$CERTS_DIR"/server.pem

# Generate the private key for the webhook server
#openssl genrsa -out runlocal-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
#openssl req -new -key runlocal-tls.key -subj "/CN=host.docker.internal" -config csr.conf \
#    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -extfile csr.conf -extensions req_ext -out runlocal-tls.crt
