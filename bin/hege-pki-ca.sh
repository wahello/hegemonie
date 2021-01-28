#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

THIS_FILE=$(realpath "${BASH_SOURCE[0]}")
THIS_DIR=$(dirname "${THIS_FILE}")
REPO=$(readlink -e "${THIS_DIR}/..")

[[ -n "$1" ]]
[[ -d "$1" ]]
D="$1"

cat >>"$D/certificate.conf" <<EOF
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
DNS.1 = localhost
DNS.2 = localhost.localdomain
DNS.3 = hege_regions
DNS.4 = hege_events
DNS.5 = hege_maps
IP.1 = 127.0.0.1
EOF

openssl genrsa \
	-out "$D/ca.key" 4096

openssl req \
	-new -x509 \
	-key "$D/ca.key" \
	-sha256 \
	-subj "/C=FR/ST=Nord/O=Hegemonie/CN=localhost" \
	-days 365 \
	-out "$D/ca.cert"

chmod 0444 "$D/ca.cert" "$D/ca.key"

