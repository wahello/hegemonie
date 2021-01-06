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

cd "$D"


#-----------------------------------------------------------------------------#

mkdir maps
$REPO/bin/hege-map-transform.sh $REPO/docs/maps/ ./maps
hege server map --endpoint=localhost:8083 ./maps &


mkdir events
hege server event --endpoint=localhost:8082 ./events &


mkdir ./live
cp -rp $REPO/docs/definitions/hegeIV ./defs
hege server region --endpoint=localhost:8081 ./defs ./live &


#-----------------------------------------------------------------------------#

trap finish SIGTERM SIGINT
wait

