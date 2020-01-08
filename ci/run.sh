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
	-north 127.0.0.1:8080 \
	-region 127.0.0.1:8081 \
	&

hegemonie region \
	-load $CONFIG \
	-save /tmp \
	-north 127.0.0.1:8081 \
	&

trap finish SIGTERM SIGINT
wait
