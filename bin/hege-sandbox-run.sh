#!/bin/bash
set -euxo pipefail

THIS_FILE=$(realpath "${BASH_SOURCE[0]}")
THIS_DIR=$(dirname "${THIS_FILE}")
REPO=$(readlink -e "${THIS_DIR}/..")

D=$(mktemp -d)

function finish() {
	set +e
	kill %3 %2 %1
	wait
}

cd "$D"


#-----------------------------------------------------------------------------#

mkdir maps
$REPO/bin/hege-map-transform.sh $REPO/docs/maps/ ./maps
hege server map --endpoint=localhost:8083 ./maps &


mkdir events
hege server event --endpoint=localhost:8082 ./events &


mkdir ./live
cp -rp $REPO/docs/definitions/hegeIV ./defs
hege server region --endpoint=localhost:8081 ./defs ./live &


#-----------------------------------------------------------------------------#

trap finish SIGTERM SIGINT
wait

