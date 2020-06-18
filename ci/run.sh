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
#   $ run.sh MAP DEFINITIONS TRANSLATIONS
# with:
#   MAP:          the path to the map seed, a.k.a a JSON file containing an object that coarsely describe the map and the name of the cities
#   DEFINITIONS:  the path to the definitions that make the world (knowledge, buildings, troops)
#   TRANSLATIONS: the path to the directory with all the go-i18n TOML files
#
# Example:
#   ./ci/run.sh \
#     ./docs/hegeIV/map-calaquyr.json \
#     ./docs/hegeIV/definitions \
#     ./docs/hegeIV/lang \
#

set -ex

# Sanitize the input
BASE=$1
[[ -d "${BASE}" ]]
[[ -d "${BASE}/definitions" ]]
[[ -d "${BASE}/lang" ]]
[[ -d "${BASE}/live" ]]
[[ -d "${BASE}/save" ]]
[[ -d "${BASE}/evt" ]]
shift

DEFS="${BASE}/definitions"
[[ -r "${DEFS}/config.json" ]]
[[ -r "${DEFS}/units.json" ]]
[[ -r "${DEFS}/buildings.json" ]]
[[ -r "${DEFS}/knowledge.json" ]]

TRANSLATIONS="${BASE}/lang"
[[ -d "${TRANSLATIONS}" ]]
[[ -r "${TRANSLATIONS}/active.en.toml" ]]

# Spawn the core services
function finish() {
	set +e
	kill %4 %3 %2 %1
	wait
}

heged auth \
	--id hege,aaa,1 \
	--live "${BASE}/live" \
	--save "${BASE}/save" \
	--endpoint 127.0.0.1:8082 \
	&

heged evt \
	--id hege,evt,1 \
	--base "${BASE}/evt" \
	--endpoint 127.0.0.1:8083 \
	&

heged region \
	--id hege,reg,1 \
	--defs "${BASE}/definitions" \
	--live "${BASE}/live" \
	--save "${BASE}/save" \
	--endpoint 127.0.0.1:8081 \
	--event 127.0.0.1:8083 \
	&

heged web \
	--id hege,web,1 \
	--lang "${BASE}/lang" \
	--templates "${PWD}/pkg/web/templates" \
	--static "${PWD}/pkg/web/static" \
	--endpoint 127.0.0.1:8080 \
	--region 127.0.0.1:8081 \
	--auth 127.0.0.1:8082 \
	--event 127.0.0.1:8083 \
	&

trap finish SIGTERM SIGINT
wait
