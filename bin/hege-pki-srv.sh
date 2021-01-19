#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

THIS_FILE=$(realpath "${BASH_SOURCE[0]}")
THIS_DIR=$(dirname "${THIS_FILE}")
REPO=$(readlink -e "${THIS_DIR}/..")

D="$1"
[[ -n "$D" ]] && [[ -d "$D" ]]
[[ -r "$D/certificate.conf" ]]
[[ -r "$D/ca.cert" ]]
[[ -r "$D/ca.key" ]]

srv="$2"
[[ -n "$srv" ]] && [[ "$srv" != "ca" ]]

openssl genrsa \
	-out "$D/$srv.key" 4096

openssl req \
	-new \
	-key "$D/$srv.key" \
	-config "$D/certificate.conf" \
	-out "$D/$srv.csr"

openssl x509 \
	-req \
	-in "$D/$srv.csr" \
	-CA "$D/ca.cert" \
	-CAkey "$D/ca.key" \
	-CAcreateserial \
	-out "$D/$srv.crt" \
	-days 365 \
	-sha256 \
	-extfile "$D/certificate.conf" \
	-extensions req_ext

rm "$D/$srv.csr"
cat "$D/$srv.crt" "$D/$srv.key" > "$D/$srv.pem"
chmod 0444 "$D/$srv.key" "$D/$srv.crt" "$D/$srv.pem"
echo "$D/$srv.pem" "$D/$srv.key"
