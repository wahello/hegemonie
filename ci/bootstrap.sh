#!/usr/bin/env bash

# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

#
# Description:
#   bootstrap.sh generates a complete configuration enough to run a minimal environment,
#   made of a single map region.
#
# Usage:
#   $ bootstrap.sh MAP DEFINITIONS OUTPUT
# with:
#   MAP:         the path to the map seed, a.k.a a JSON file containing an object that coarsely describe the map and the name of the cities
#   DEFINITIONS: the path to the definitions that make the world (knowledge, buildings, troops)
#   OUTPUT:      the path to the live DB of the world (actual map, users, cities, armies) plus a copy of the expanded description.
#

set -xe

function usage() { echo -e "USAGE:\n  $0  /path/to/map_seed.json  /path/to/definitions  /path/to/live" ; }
function error() { echo $@ 1>&2 ; exit 2 ; }
function check_input() { [ -r "$1" ] || ( usage ; error "Missing $1" ); }
function check_file() { [ -r "$1" ] || ( error "Missing $1" ); }


# Sanitize the input
MAP=$1
check_input "${MAP}"
shift

DEFS=$1
[ -d "$DEFS" ] || error "Missing DEFINITIONS folder"
check_input "${DEFS}/config.json"
check_input "${DEFS}/units.json"
check_input "${DEFS}/buildings.json"
check_input "${DEFS}/knowledge.json"
shift

OUT=$1
[[ -d "${OUT}" ]] || error "Invalid output path: [$OUT]"
shift
mkdir -p $OUT/definitions
mkdir -p $OUT/live
mkdir -p $OUT/save

TMP=$(mktemp -d)


# Generate the database for a world
hege-mapper normalize < "${MAP}" > "${TMP}/map_seed.json"
hege-mapper export --config "${DEFS}" "${TMP}" < "${TMP}/map_seed.json" > "${TMP}/env"
. "${TMP}/env"

check_file "${HEGE_LIVE}/cities.json"
check_file "${HEGE_LIVE}/map.json"
check_file "${HEGE_LIVE}/armies.json"
check_file "${HEGE_LIVE}/fights.json"
check_file "${TMP}/auth.json"

cp -v -p \
  "${HEGE_DEFS}/units.json" \
  "${HEGE_DEFS}/config.json" \
  "${HEGE_DEFS}/buildings.json" \
  "${HEGE_DEFS}/knowledge.json" \
	"${TMP}/map_seed.json" \
  "${OUT}/definitions/"

cp -v -p \
  "${HEGE_LIVE}/armies.json" \
  "${HEGE_LIVE}/fights.json" \
  "${HEGE_LIVE}/cities.json" \
  "${HEGE_LIVE}/map.json" \
  "${TMP}/auth.json" \
	"${OUT}/live/"

echo $OUT
