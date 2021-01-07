#!/bin/sh
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

# To enable this hook, rename this file to "pre-commit" and copy it in the
# .git/hooks directory

exec 1>&2

rc=0



#------------------------------------------------------------------------------
# Check that any code file has a copyright preamble that mentions the current
# year. Auto-generated files are not considered.

Y=$(date +%Y)

files_without_proper_copyright() {
	grep -rnIL "Copyright (c) 2018-$Y Contributors as noted in the AUTHORS file" -- * \
		| egrep -e '.(go|py|sh)$' \
		| egrep -v -e '.pb.(go|py|sh)'
}

if [ $(files_without_proper_copyright | wc -l) -gt 0 ] ; then
	echo "### Files lacking a proper copyright preamble:"
	echo "# Copyright (c) 2018-$Y Contributors as noted in the AUTHORS file"
	echo "# This Source Code Form is subject to the terms of the Mozilla Public"
	echo "# License, v. 2.0. If a copy of the MPL was not distributed with this"
	echo "# file, You can obtain one at http://mozilla.org/MPL/2.0/."
	files_without_proper_copyright
	rc=1
fi



#------------------------------------------------------------------------------
# Check that any modified code is well formated

if [ $(gofmt -l . | wc -l) -gt 0 ] ; then
	echo "### Files lacking a proper format:"
	gofmt -l .
	rc=1
fi

exit $rc

