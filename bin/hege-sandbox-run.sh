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
	kill %4 %3 %2 %1
	wait
}

trap finish SIGTERM SIGINT EXIT

cd "$D"


#-----------------------------------------------------------------------------
# Generate a self-signed certificate
mkdir pki
$REPO/bin/hege-pki-ca.sh $D/pki


function tls_opts() {
	local srv base
	srv=$1 ; shift
	base=$D/pki
	echo "--tls-key=$base/$srv.key" "--tls-crt=$base/$srv.crt"
}


#-----------------------------------------------------------------------------
#
mkdir maps
"$REPO/bin/hege-map-transform.sh" "$REPO/docs/maps/" "$D/maps"
"$REPO/bin/hege-pki-srv.sh" "$D/pki" maps
hege server $(tls_opts maps) --endpoint=localhost:8083 maps --defs="$D/maps" &


#-----------------------------------------------------------------------------
# 
mkdir events
"$REPO/bin/hege-pki-srv.sh" "$D/pki" events
hege server $(tls_opts events) --endpoint=localhost:8082 events "$D/events" &


#-----------------------------------------------------------------------------
# 
mkdir live 
cp -rp "$REPO/docs/definitions/hegeIV" "$D/defs"
"$REPO/bin/hege-pki-srv.sh" "$D/pki" regions
hege server $(tls_opts regions) --endpoint=localhost:8081 regions --defs="$D/defs" "$D/live" &


#-----------------------------------------------------------------------------
# Generate a per-service certificate
mkdir proxy
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

frontend public
	bind :8000  alpn http/1.1
	bind :8443  ssl  verify none  crt $D/pki/proxy.pem alpn h2,http/1.1
	http-request deny unless { req.hdr(content-type) -m str "application/grpc" }
	acl ismap   path_beg "/hege.map."
	acl isevt   path_beg "/hege.evt."
	acl isreg   path_beg "/hege.reg."
	use_backend grpc_map    if ismap
	use_backend grpc_evt    if isevt
	use_backend grpc_reg    if isreg

backend grpc_reg
	balance roundrobin
	server reg1 localhost:8081  ssl  alpn h2  maxconn 32  verify none

backend grpc_evt
	balance roundrobin
	server evt1 localhost:8082  ssl  alpn h2  maxconn 32  verify none

backend grpc_map
	balance roundrobin
	server map1 localhost:8083  ssl  alpn h2  maxconn 32  verify none

EOF

"$REPO/bin/hege-pki-srv.sh" "$D/pki" proxy
haproxy -- "$D/proxy/haproxy.cfg" &


wait

