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

L='(C) Quentin Minten / CC BY-NC-SA 3.0'

for F in "$SRC/quentin-minten*/*.jpg" ; do
  BN=$(basename $F)
  convert "$F" \
    -gravity south \
    -stroke none \
    -fill yellow \
    -strokewidth 2 \
    -annotate 0 "$L" \
    "$DST/$BN"
done
