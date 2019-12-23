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

hege-front \
	-templates $PWD/hege-front/templates \
	-static $PWD/hege-front/static \
	&

hege-world \
	-load $CONFIG \
	&

trap finish SIGTERM SIGINT
wait
