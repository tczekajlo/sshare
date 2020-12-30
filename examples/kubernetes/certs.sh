#!/bin/bash

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=DE/ST=/L=/O=/OU=/CN=sshare"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate private key and certificate
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=DE/ST=/L=/O=/OU=/CN=sshare" -config san.cnf

# 3. Use CA's private key to sign CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 600 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extensions req_ext -extfile san.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text
