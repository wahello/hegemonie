#!/usr/bin/env bash

# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

#
# Description:
#   bootstrap.sh generates a complete configuration enough to run a minimal environment,
#   made of a single region whose map is given as an argument.
#
# Usage:
#   $ bootstrap.sh MAP DEFINITIONS OUTPUT
# with:
#   MAP:         the path to the map seed, a.k.a a JSON file containing an object that coarsely describe the map and the name of the cities
#   DEFINITIONS: the path to the definitions that make the world (knowledge, buildings, troops)
#   OUTPUT:      the path to the live DB of the world (actual map, users, cities) plus a copy of the expanded description.
#
# Example:
#   ./ci/bootstrap.sh \
#     ./docs/hegeIV/map-calaquyr.json \
#     ./docs/hegeIV/definitions \
#     $(mktemp -d)

set -ex


# Sanitize the input
MAP=$1
[[ -r "${MAP}" ]]
shift

DEFS=$1
[[ -d "${DEFS}" ]]
[[ -r "${DEFS}/config.json" ]]
[[ -r "${DEFS}/units.json" ]]
[[ -r "${DEFS}/buildings.json" ]]
[[ -r "${DEFS}/knowledge.json" ]]
shift

TRANSLATIONS=$1
[[ -d "${TRANSLATIONS}" ]]
[[ -r "${TRANSLATIONS}/active.en.toml" ]]
shift

OUT=$1
[[ -d "${OUT}" ]]
shift

mkdir -p "${OUT}/definitions"
mkdir -p "${OUT}/lang"
mkdir -p "${OUT}/live"
mkdir -p "${OUT}/save"
mkdir -p "${OUT}/evt"
TMP=$(mktemp -d)


# Generate the database for a world
hege-mapper normalize < "${MAP}" > "${TMP}/map_seed.json"
hege-mapper export --config "${DEFS}" "${TMP}" < "${TMP}/map_seed.json" > "${TMP}/env"
. "${TMP}/env"

[[ -r "${HEGE_LIVE}/cities.json" ]]
[[ -r "${HEGE_LIVE}/map.json" ]]
[[ -r "${HEGE_LIVE}/fights.json" ]]
[[ -r "${TMP}/auth.json" ]]

cp -p \
  "${HEGE_DEFS}/units.json" \
  "${HEGE_DEFS}/config.json" \
  "${HEGE_DEFS}/buildings.json" \
  "${HEGE_DEFS}/knowledge.json" \
  "${OUT}/definitions/"

cp -p \
  "${TRANSLATIONS}/"active.*.toml \
  "${OUT}/lang"

cp -p \
	"${TMP}/map_seed.json" \
  "${OUT}/"

cp -p \
  "${HEGE_LIVE}/fights.json" \
  "${HEGE_LIVE}/cities.json" \
  "${HEGE_LIVE}/map.json" \
  "${TMP}/auth.json" \
	"${OUT}/live/"

rm -rf "${TMP}"

echo "${OUT}"
