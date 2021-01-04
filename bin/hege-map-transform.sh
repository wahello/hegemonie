#!/bin/bash
set -euxo pipefail

SRC=$1 ; shift
DST=$1 ; shift
[[ -n "$SRC" ]]
[[ -n "$DST" ]]

for S in "$SRC"/*.seed.json ; do
	D=${S/.seed.json/.final.json}
	hege tools map init < "$S" | hege tools map normalize > "$D"
done

