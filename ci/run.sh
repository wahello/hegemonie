#!/usr/bin/env bash
# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

#
# ./ci/bootstrap.sh \
#   docs/hegeIV/map-calaquyr.json \
#   docs/hegeIV/definitions/ \
#   ci/hegeIV-calaquyr/
#

set -xe

function usage() { echo -e "USAGE:\n  $0  DIR_DEFS  DIR_LIVE [DIR_SAVE}" ; }
function error() { ( echo $@ ; usage ) 1>&2 ; exit 2 ; }
function check_file() { [ -r "$1" ] || error "Missing $1"; }


# Sanitize the input
DEFS=$1
[ -d "$DEFS" ] || error "Missing DEFINITIONS folder"
check_file "${DEFS}/config.json"
check_file "${DEFS}/units.json"
check_file "${DEFS}/buildings.json"
check_file "${DEFS}/knowledge.json"
shift

LIVE=$1
[ -d "${LIVE}" ] || error "Missing LIVE folder"
check_file "${LIVE}/auth.json"
check_file "${LIVE}/armies.json"
check_file "${LIVE}/cities.json"
check_file "${LIVE}/map.json"
shift

# Prepare the working environment
TMP=$(mktemp -d)
mkdir $TMP/save


# Spawn the core services
function finish() {
	set +e
	kill %3
	kill %2
	kill %1
	wait
}

hegemonie web agent \
	--templates $PWD/pkg/web/templates \
	--static $PWD/pkg/web/static \
	--endpoint 127.0.0.1:8080 \
	--region 127.0.0.1:8081 \
	--auth 127.0.0.1:8082 \
	&

hegemonie region agent \
	--defs "${DEFS}" \
	--live "${LIVE}" \
	--save "${TMP}/save" \
	--endpoint 127.0.0.1:8081 \
	&

hegemonie auth agent \
	--live "${LIVE}" \
	--save "${TMP}/save" \
	--endpoint 127.0.0.1:8082 \
	&

trap finish SIGTERM SIGINT
wait
