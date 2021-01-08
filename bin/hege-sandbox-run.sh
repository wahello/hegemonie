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


#-----------------------------------------------------------------------------#

mkdir maps
"$REPO/bin/hege-map-transform.sh" "$REPO/docs/maps/" ./maps
hege server map --endpoint=localhost:8083 ./maps &


mkdir ./events
hege server event --endpoint=localhost:8082 ./events &


mkdir ./live
cp -rp "$REPO/docs/definitions/hegeIV" ./defs
hege server region --endpoint=localhost:8081 ./defs ./live &

mkdir ./proxy

cat >>proxy/certificate.conf <<EOF
[ req ]
prompt = no
default_bits = 4096
distinguished_name = req_distinguished_name
req_extensions = req_ext
[ req_distinguished_name ]
C=FR
ST=Nord
L=Hem
O=OpenIO
OU=R&D
CN=localhost
[ req_ext ]
subjectAltName = @alt_names
[alt_names]
DNS.1 = hostname.domain.tld
DNS.2 = hostname
IP.1 = 127.0.0.1
EOF
openssl genrsa \
	-out "$D/proxy/ca.key" 4096
openssl req \
	-new -x509 \
	-key "$D/proxy/ca.key" \
	-sha256 \
	-subj "/C=FR/ST=Nord/O=CA, Inc./CN=localhost" \
	-days 365 \
	-out "$D/proxy/ca.cert"
openssl genrsa \
	-out "$D/proxy/service.key" 4096
openssl req \
	-new \
	-key "$D/proxy/service.key" \
	-config "$D/proxy/certificate.conf" \
	-out "$D/proxy/service.csr"
openssl x509 \
	-req \
	-in "$D/proxy/service.csr" \
	-CA "$D/proxy/ca.cert" \
	-CAkey "$D/proxy/ca.key" \
	-CAcreateserial \
	-out "$D/proxy/service.pem" \
	-days 365 \
	-sha256 \
	-extfile "$D/proxy/certificate.conf" \
	-extensions req_ext

# debug
openssl x509 \
	-in "$D/proxy/service.pem" \
	-text \
	-noout

( cd "$D/proxy" && chmod 400 -- * )

cat >>proxy/haproxy.cfg <<EOF
global
	log stdout local0
	stats socket $D/proxy/admin.sock mode 660 level admin expose-fd listeners
	stats timeout 30s
	# Default SSL material locations
	ca-base /etc/ssl/certs
	crt-base /etc/ssl/private
	# See: https://ssl-config.mozilla.org/#server=haproxy&server-version=2.0.3&config=intermediate
	ssl-default-bind-ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384
	ssl-default-bind-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256
	ssl-default-bind-options ssl-min-ver TLSv1.2 no-tls-tickets

defaults
	log	global
	mode http
	option httplog
	option logasap
	option http-use-htx
	timeout connect 2s
	timeout client  30s
	timeout server  30s

frontend grpc
	#bind :8080  alpn h2,h2c
	bind 127.0.0.1:8080  ssl  verify none  crt $D/proxy/service.pem alpn h2
	http-request deny unless { req.hdr(content-type) -m str "application/grpc" }
	acl ismap  path_beg "/hege.map."
	acl isevt  path_beg "/hege.evt."
	acl isreg  path_beg "/hege.reg."
	use_backend grpc_map  if ismap
	use_backend grpc_evt  if isevt
	use_backend grpc_reg  if isreg

backend grpc_map
	balance roundrobin
	server map1 localhost:8081  ssl verify none  alpn h2  check

backend grpc_evt
	balance roundrobin
	server evt1 localhost:8082  ssl verify none  alpn h2  check

backend grpc_reg
	balance roundrobin
	server reg1 localhost:8083  ssl verify none  alpn h2  check

EOF

haproxy -d -V -- "$D/proxy/haproxy.cfg"

#-----------------------------------------------------------------------------#

wait

