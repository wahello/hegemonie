// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

func (r *Resources) GreaterOrEqualTo(o *Resources) bool {
	for i := 0; i < ResourceMax; i++ {
		if r[i] < o[i] {
			return false
		}
	}
	return true
}

func (r *Resources) Add(o *Resources) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = r[i] + o[i]
	}
}

func (r *Resources) Remove(o *Resources) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = r[i] - o[i]
	}
}

func (r *Resources) TrimTo(limit *Resources) {
	for i := 0; i < ResourceMax; i++ {
		if r[i] < limit[i] {
			r[i] = r[i]
		} else {
			r[i] = limit[i]
		}
	}
}
