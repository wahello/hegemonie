#!/bin/bash
set -euxo pipefail

D="$1" ; shift
[[ -n "$D" ]]

for SRC in $D/*.seed.json ; do
	DST=${SRC/.seed./.final.}
	hege tools map init < "$SRC" | hege tools map normalize > "$DST"
done

