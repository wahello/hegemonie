#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

SRC=$1 ; shift
DST=$1 ; shift
[[ -n "$SRC" ]]
[[ -n "$DST" ]]

for S in "$SRC"/*.seed.json ; do
	D=$(basename "$S")
	D=${D/.seed.json/.final.json}
	hege tools map init < "$S" | hege tools map normalize > "$DST/$D"
done

