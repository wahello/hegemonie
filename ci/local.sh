#!/usr/bin/env bash
# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -ex
make -s
DST=$(mktemp -d)

nice ionice ./ci/bootstrap.sh \
  docs/hegeIV/map-calaquyr.json \
  docs/hegeIV/definitions \
  docs/hegeIV/lang \
  $DST

find "$DST" -type f

nice ionice ./ci/run.sh $DST
