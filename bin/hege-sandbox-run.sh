#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

THIS_FILE=$(realpath "${BASH_SOURCE[0]}")
THIS_DIR=$(dirname "${THIS_FILE}")
REPO=$(readlink -e "${THIS_DIR}/..")

D=$(mktemp -d)

function finish() {
	set +e
	kill %3 %2 %1
	wait
}

trap finish SIGTERM SIGINT EXIT

cd "$D"
mkdir maps events live proxy
mkdir pki

#-----------------------------------------------------------------------------
# Generate a self-signed certificate
cat >>"$D/pki/certificate.conf" <<EOF
[ req ]
prompt = no
default_bits = 4096
distinguished_name = req_distinguished_name
req_extensions = req_ext
[ req_distinguished_name ]
C=FR
ST=Nord
L=Flines-lez-Mortagne
O=Hegemonie
OU=R&D
CN=localhost
[ req_ext ]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost.localdomain
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF
openssl genrsa \
	-out "$D/pki/ca.key" 4096
openssl req \
	-new -x509 \
	-key "$D/pki/ca.key" \
	-sha256 \
	-subj "/C=FR/ST=Nord/O=Hegemonie/CN=localhost" \
	-days 365 \
	-out "$D/pki/ca.cert"

function service_certificate() {
	local srv
	srv=$1 ; shift
	openssl genrsa \
		-out "$D/pki/$srv.key" 4096
	openssl req \
		-new \
		-key "$D/pki/$srv.key" \
		-config "$D/pki/certificate.conf" \
		-out "$D/pki/$srv.csr"
	openssl x509 \
		-req \
		-in "$D/pki/$srv.csr" \
		-CA "$D/pki/ca.cert" \
		-CAkey "$D/pki/ca.key" \
		-CAcreateserial \
		-out "$D/pki/$srv.crt" \
		-days 365 \
		-sha256 \
		-extfile "$D/pki/certificate.conf" \
		-extensions req_ext
	rm "$D/pki/$srv.csr"
	cat "$D/pki/$srv.crt" "$D/pki/$srv.key" > "$D/pki/$srv.pem"
	chmod 0400 "$D/pki/$srv.key" "$D/pki/$srv.crt"
	echo "$D/pki/$srv.pem" "$D/pki/$srv.key"
}

function tls_opts() {
	local srv
	srv=$1 ; shift
	echo "--key=$D/pki/$srv.key" "--crt=$D/pki/$srv.crt"
}

#-----------------------------------------------------------------------------
#
"$REPO/bin/hege-map-transform.sh" "$REPO/docs/maps/" "$D/maps"
service_certificate maps
hege server $(tls_opts maps) map --endpoint=localhost:8083 "$D/maps" &


#-----------------------------------------------------------------------------
# 
service_certificate events
hege server $(tls_opts events) event --endpoint=localhost:8082 "$D/events" &


#-----------------------------------------------------------------------------
# 
cp -rp "$REPO/docs/definitions/hegeIV" "$D/defs"
service_certificate region
hege server $(tls_opts region) region --endpoint=localhost:8081 "$D/defs" "$D/live" &


#-----------------------------------------------------------------------------
# Generate a per-service certificate


cat >>"$D/proxy/haproxy.cfg" <<EOF
global
	log stdout local0
	stats socket $D/pki/admin.sock mode 660 level admin expose-fd listeners
	stats timeout 30s
	ca-base /etc/ssl/certs
	crt-base /etc/ssl/private
	# See: https://ssl-config.mozilla.org/#server=haproxy&server-version=2.0.3&config=intermediate
	#ssl-default-bind-ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384
	#ssl-default-bind-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256
	#ssl-default-bind-options ssl-min-ver TLSv1.2 no-tls-tickets
	debug

defaults
	log global
	maxconn 1024
	mode http
	timeout connect 10s
	timeout client  30s
	timeout server  30s
	option logasap
	option httplog
	option http-use-htx

frontend grpc
	bind :8080  ssl  verify none  crt $D/pki/proxy.pem alpn h2
	http-request deny unless { req.hdr(content-type) -m str "application/grpc" }
	acl ismap  path_beg "/hege.map."
	acl isevt  path_beg "/hege.evt."
	acl isreg  path_beg "/hege.reg."
	use_backend grpc_map  if ismap
	use_backend grpc_evt  if isevt
	use_backend grpc_reg  if isreg

backend grpc_reg
	balance roundrobin
	server reg1 localhost:8081  ssl  alpn h2  check  maxconn 32  verify none

backend grpc_evt
	balance roundrobin
	server evt1 localhost:8082  ssl  alpn h2  check  maxconn 32  verify none

backend grpc_map
	balance roundrobin
	server map1 localhost:8083  ssl  alpn h2  check  maxconn 32  verify none

EOF

service_certificate proxy
haproxy -- "$D/proxy/haproxy.cfg"

#-----------------------------------------------------------------------------#

wait

