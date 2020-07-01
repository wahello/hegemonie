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


# Generate the database for a world
hege-mapper normalize < "${MAP}" > "${OUT}/map_seed.json"
hege-mapper export --config "${DEFS}" "${OUT}" < "${OUT}/map_seed.json"

[[ -r "${OUT}/definitions/config.json" ]]
[[ -r "${OUT}/definitions/knowledge.json" ]]
[[ -r "${OUT}/definitions/buildings.json" ]]
[[ -r "${OUT}/definitions/units.json" ]]
[[ -r "${OUT}/live/map.json" ]]
[[ -r "${OUT}/live/cities.json" ]]
[[ -r "${OUT}/live/fights.json" ]]
[[ -r "${OUT}/auth.json" ]]

cp -p \
  "${TRANSLATIONS}/"active.*.toml \
  "${OUT}/lang"

echo "${OUT}"
