#!/usr/bin/env bash
## -*- coding: utf-8 -*-
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail
THIS_FILE=$(realpath "${BASH_SOURCE[0]}")
THIS_DIR=$(dirname "$THIS_FILE")
BASE_DIR=$(dirname "$THIS_DIR")
pushd "${BASE_DIR}"

case $1 in
	env)
		virtualenv .env
		. .env/bin/activate
		pip3 install --upgrade -r requirements.txt
		;;
	www)
		. .env/bin/activate
		./bin/wwwgen.py src .build
		;;
	*)
		echo "$0 (www|env)" 1>&2
		exit 2
		;;
esac

