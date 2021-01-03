#!/bin/bash
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
