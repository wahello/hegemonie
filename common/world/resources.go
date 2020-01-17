// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

func (r Resources) Equals(o Resources) bool {
	for i := 0; i < ResourceMax; i++ {
		if r[i] != o[i] {
			return false
		}
	}
	return true
}

func (r Resources) GreaterOrEqualTo(o Resources) bool {
	for i := 0; i < ResourceMax; i++ {
		if r[i] < o[i] {
			return false
		}
	}
	return true
}

func (r Resources) GreaterThan(o Resources) bool {
	for i := 0; i < ResourceMax; i++ {
		if r[i] <= o[i] {
			return false
		}
	}
	return true
}

func (r Resources) IsZero() bool {
	for i := 0; i < ResourceMax; i++ {
		if r[i] != 0 {
			return false
		}
	}
	return true
}

func (r Resources) GetRatio(nb float64) Resources {
	var rc Resources = r
	for i := 0; i < ResourceMax; i++ {
		rc[i] = uint64(float64(rc[i]) * nb)
	}
	return rc
}

func (r *Resources) Zero() {
	for i := 0; i < ResourceMax; i++ {
		r[i] = 0
	}
}

func (r *Resources) Add(o Resources) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = r[i] + o[i]
	}
}

func (r *Resources) Remove(o Resources) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = r[i] - o[i]
	}
}

func (r *Resources) TrimTo(limit Resources) {
	for i := 0; i < ResourceMax; i++ {
		if r[i] < limit[i] {
			r[i] = r[i]
		} else {
			r[i] = limit[i]
		}
	}
}

func (r *Resources) Multiply(m ResourcesMultiplier) {
	for i := 0; i < ResourceMax; i++ {
		vf := float64(r[i]) * m[i]
		if vf < 0 {
			r[i] = 0
		} else {
			r[i] = uint64(vf)
		}
	}
}

func MultiplierUniform(nb float64) ResourcesMultiplier {
	var rc ResourcesMultiplier
	for i := 0; i < ResourceMax; i++ {
		rc[i] = nb
	}
	return rc
}
