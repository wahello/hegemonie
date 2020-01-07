#!/usr/bin/env bash
set -x
set -e

CONFIG=$1 ; shift

function finish() {
	set +e
	kill %2
	kill %1
	wait
}

hegemonie front \
	-templates $PWD/front/templates \
	-static $PWD/front/static \
	&

hegemonie region \
	-load $CONFIG \
	&

trap finish SIGTERM SIGINT
wait
