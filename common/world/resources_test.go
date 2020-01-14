// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"testing"
)

func TestResources(t *testing.T) {
	zero := Resources{}
	r0 := Resources{}
	r1 := Resources{}
	t.Log(r0)
	t.Log(r1)
	if !zero.IsZero() || !r0.IsZero() || !r1.IsZero() {
		t.Fatal()
	}
	if !r0.GreaterOrEqualTo(r0) {
		t.Fatal()
	}
	r1[1] = 1
	if !r1.GreaterOrEqualTo(r1) {
		t.Fatal()
	}
	if !r1.GreaterOrEqualTo(r0) || r0.GreaterOrEqualTo(r1) {
		t.Fatal()
	}

	r0.Add(r1)
	if !r0.Equals(r1) {
		t.Fatal()
	}
	r0.TrimTo(zero)
	if !r0.IsZero() {
		t.Fatal()
	}
}
