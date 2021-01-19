#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

function clean() { docker image remove "jfsmig/hegemonie-$1" ; }

function build() { docker build --target="$1" --tag="jfsmig/hegemonie-$1" . ; }

function foreach() { for T in dependencies runtime demo debug demo-prometheus ; do $1 $T ; done ; }

function pushall() { for T in runtime demo demo-prometheus ; do docker push "jfsmig/hegemonie-$T:latest" ; done ; }

if [[ $# == 0 ]] ; then
	foreach build
	pushall
else
	case $1 in
			build) foreach build ;;
			push) pushall ;;
			clean) foreach clean ;;
	esac
fi

