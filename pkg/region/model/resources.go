// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

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

func (r *Resources) Set(o Resources) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = o[i]
	}
}

func (r *Resources) SetValue(v uint64) {
	for i := 0; i < ResourceMax; i++ {
		r[i] = v
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
		if r[i] > limit[i] {
			r[i] = limit[i]
		}
	}
}

func (r *Resources) Multiply(rm ResourcesMultiplier) {
	for i := 0; i < ResourceMax; i++ {
		vf := float64(r[i]) * rm[i]
		if vf < 0 {
			r[i] = 0
		} else {
			r[i] = uint64(vf)
		}
	}
}

func (r *Resources) Increment(ri ResourcesIncrement) {
	for i := 0; i < ResourceMax; i++ {
		post := r[i] + uint64(ri[i])
		if post > r[i] {
			r[i] = 0
		} else {
			r[i] = post
		}
	}
}

func (r *Resources) Apply(rm ResourceModifiers) {
	r.Multiply(rm.Mult)
	r.Increment(rm.Plus)
}

func (rm *ResourcesMultiplier) SetValue(v float64) {
	for i := 0; i < ResourceMax; i++ {
		rm[i] = v
	}
}

func (ri *ResourcesIncrement) SetValue(v int64) {
	for i := 0; i < ResourceMax; i++ {
		ri[i] = v
	}
}

func MultiplierUniform(nb float64) (rc ResourcesMultiplier) {
	rc.SetValue(nb)
	return rc
}

func IncrementUniform(nb int64) (rc ResourcesIncrement) {
	rc.SetValue(nb)
	return rc
}

func ResourcesUniform(nb uint64) (rc Resources) {
	rc.SetValue(nb)
	return rc
}

func ResourceModifierUniform(mult float64, inc int64) ResourceModifiers {
	return ResourceModifiers{
		MultiplierUniform(mult),
		IncrementUniform(inc),
	}
}

func ResourceModifierNoop() ResourceModifiers {
	return ResourceModifierUniform(1.0, 0.0)
}

func (o0 *ResourceModifiers) ComposeWith(o1 ResourceModifiers) {
	for i := 0; i < ResourceMax; i++ {
		o0.Mult[i] *= o1.Mult[i]
		o0.Plus[i] += o1.Plus[i]
	}
}
