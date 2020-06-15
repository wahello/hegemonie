#!/usr/bin/env bash
# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

#
# Description:
#   run.sh generates a complete configuration enough to run a minimal environment,
#   and start one daemon of each kind. Each daemon will be bond to its default port
#   on 127.0.0.1
#
# Usage:
#   $ run.sh MAP DEFINITIONS
# with:
#   MAP:         the path to the map seed, a.k.a a JSON file containing an object that coarsely describe the map and the name of the cities
#   DEFINITIONS: the path to the definitions that make the world (knowledge, buildings, troops)
#
# Example:
#   ./ci/run.sh \
#     ./docs/hegeIV/map-calaquyr.json \
#     ./docs/hegeIV/definitions \
#

set -ex

# Sanitize the input
MAP=$1
[[ -r "${MAP}" ]]
shift

DEFS=$1
[[ -d "$DEFS" ]]
[[ "${DEFS}/config.json" ]]
[[ "${DEFS}/units.json" ]]
[[ "${DEFS}/buildings.json" ]]
[[ "${DEFS}/knowledge.json" ]]
shift

# Prepare the working environment
TMP=$(mktemp -d)
mkdir $TMP/live
mkdir $TMP/save
mkdir $TMP/evt

./ci/bootstrap.sh "${MAP}" "${DEFS}" "${TMP}"


# Spawn the core services
function finish() {
	set +e
	kill %4 %3 %2 %1
	wait
}

heged auth \
	--id hege,aaa,1 \
	--live "${TMP}/live" \
	--save "${TMP}/save" \
	--endpoint 127.0.0.1:8082 \
	&

heged evt \
	--id hege,evt,1 \
	--base "${TMP}/evt" \
	--endpoint 127.0.0.1:8083 \
	&

heged region \
	--id hege,reg,1 \
	--defs "${DEFS}" \
	--live "${TMP}/live" \
	--save "${TMP}/save" \
	--endpoint 127.0.0.1:8081 \
	--event 127.0.0.1:8083 \
	&

heged web \
	--id hege,web,1 \
	--templates $PWD/pkg/web/templates \
	--static $PWD/pkg/web/static \
	--endpoint 127.0.0.1:8080 \
	--region 127.0.0.1:8081 \
	--auth 127.0.0.1:8082 \
	--event 127.0.0.1:8083 \
	&

trap finish SIGTERM SIGINT
wait
