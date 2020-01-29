#!/usr/bin/env bash
set -x
set -e

CONFIG=$1 ; shift

function finish() {
	set +e
	kill %3
	kill %2
	kill %1
	wait
}

hegemonie web agent \
	--templates $PWD/pkg/web/templates \
	--static $PWD/pkg/web/static \
	--endpoint 127.0.0.1:8080 \
	--region 127.0.0.1:8081 \
	--auth 127.0.0.1:8082 \
	&

hegemonie region agent \
	--load $CONFIG \
	--save /tmp \
	--endpoint 127.0.0.1:8081 \
	&

hegemonie auth agent \
	--load $CONFIG/auth.json \
	--save /tmp \
	--endpoint 127.0.0.1:8082 \
	&

trap finish SIGTERM SIGINT
wait
