// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMapInitFails(t *testing.T) {
	testFail := func(encoded string) {
		if err := NewMap().LoadJSON(encoded); err == nil {
			t.Fatal("Unexpected success with", encoded)
		}
	}

	// check missing or empty name
	testFail(`{"sites":[{"id":1}],"roads":[{"src":1, "dst":1}]}`)
	testFail(`{"id": "", "sites":[{"id":1}],"roads":[{"src":1, "dst":1}]}`)
	// Loop detected
	testFail(`{"id": "test map", "sites":[{"id":1}],"roads":[{"src":1, "dst":1}]}`)
	// Dangling
	testFail(`{"id": "test map", "sites":[{"id":1}],"roads":[{"src":1, "dst":2}]}`)
	testFail(`{"id": "test map", "sites":[{"id":1}],"roads":[{"src":2, "dst":1}]}`)
	// EINVAL
	testFail(`{"id": "test map", "sites":[{"id":1}],"roads":[{"src":0, "dst":1}]}`)
	testFail(`{"id": "test map", "sites":[{"id":1}],"roads":[{"src":1, "dst":0}]}`)
	testFail(`{"id": "test map", "sites":[{"id":0},{"id":1},{"id":2}],"roads":[{"src":1, "dst":2}]}`)
	// Duplications
	testFail(`{"id": "test map", "sites":[{"id":1},{"id":1},{"id":2}],"roads":[{"src":1, "dst":2},{"src":2, "dst":1}]}`)
	testFail(`{"id": "test map", "sites":[{"id":1},{"id":2}],"roads":[{"src":1, "dst":2},{"src":1, "dst":2},{"src":2, "dst":1}]}`)
}

func TestMapInitsuccess(t *testing.T) {
	testOk := func(encoded string) *Map {
		m := NewMap()
		if err := m.LoadJSON(encoded); err != nil {
			panic(fmt.Sprintf("Err:%v Map:%v", err, m))
		}
		return m
	}

	// Empty maps
	testOk(`{"id":"test"}`)
	testOk(`{"id":"test", "sites":[]}`)
	testOk(`{"id":"test", "roads":[]}`)
	testOk(`{"id":"test", "sites":[],"roads":[]}`)

	// Singleton map
	var m *Map
	v0 := rand.Uint64()
	m = testOk(fmt.Sprintf(`{"id":"test", "sites":[{"id":%v}], "roads":[]}`, v0))
	if !m.CellHas(v0) {
		t.Fatal(*m)
	}

	// Minimal map, test of roads
	v1 := rand.Uint64()
	m = testOk(fmt.Sprintf(`{"id":"test", "sites":[{"id":%v},{"id":%v}], "roads":[{"src":%v, "dst":%v}, {"src":%v, "dst":%v}]}`, v0, v1, v0, v1, v1, v0))
	if !m.CellHas(v0) {
		t.Fatal()
	}
	if !m.RoadHas(v0, v1) {
		t.Fatal()
	}
	if !m.RoadHas(v1, v0) {
		t.Fatal()
	}
}
