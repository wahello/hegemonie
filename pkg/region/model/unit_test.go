// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"sort"
	"testing"
)

type SliceTester func(marker uint64, max uint32, expectedId ...uint64)

func TestSetOfUnitType(t *testing.T) {
	u2 := &UnitType{ID: 2}
	s := SetOfUnitTypes{}
	s.Add(&UnitType{ID: 1})
	s.Add(&UnitType{ID: 3})
	s.Add(u2)
	s.Add(&UnitType{ID: 4})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get(2) {
		t.Fatal()
	}

	makeTester := func(extractor func(itf interface{}) uint64) SliceTester {
		return func(marker uint64, max uint32, expectedId ...uint64) {
			s12 := s.Slice(marker, max)
			if uint32(len(s12)) > max {
				t.Fatal()
			}
			if len(expectedId) == 0 && len(s12) > 0 {
				t.Fatal()
			}
			for idx, id := range expectedId {
				if extractor(s12[idx]) != id {
					t.Fatal()
				}
			}
		}
	}

	testSlice := makeTester(func(itf interface{}) uint64 {
		return itf.(*UnitType).ID
	})

	testSlice(0, 1, 1)
	testSlice(0, 2, 1, 2)
	testSlice(0, 3, 1, 2, 3)
	testSlice(0, 4, 1, 2, 3, 4)
	testSlice(0, 5, 1, 2, 3, 4)

	testSlice(1, 1, 2)
	testSlice(1, 2, 2, 3)
	testSlice(1, 3, 2, 3, 4)
	testSlice(1, 4, 2, 3, 4)
	testSlice(1, 5, 2, 3, 4)

	testSlice(4, 1)
}

func TestSetOfUnit(t *testing.T) {
	u2 := &Unit{ID: "2", Type: 1}
	s := SetOfUnits{}
	s.Add(&Unit{ID: "1", Type: 1})
	s.Add(&Unit{ID: "3", Type: 1})
	s.Add(u2)
	s.Add(&Unit{ID: "4", Type: 1})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get("2") {
		t.Fatal()
	}
}

func TestUnitFrontier(t *testing.T) {
	ut := SetOfUnitTypes{}
	ut.Add(&UnitType{ID: 1})
	ut.Add(&UnitType{ID: 2, RequiredBuilding: 2})
	ut.Add(&UnitType{ID: 3, RequiredBuilding: 2})
	ut.Add(&UnitType{ID: 4, RequiredBuilding: 3})

	b := SetOfBuildings{}
	b.Add(&Building{ID: "1", Type: 1})
	b.Add(&Building{ID: "2", Type: 2})
	b.Add(&Building{ID: "3", Type: 3})

	var f []*UnitType

	// Units without requirement
	f = ut.Frontier(SetOfBuildings{})
	if len(f) != 1 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}})
	if len(f) != 1 {
		t.Fatal()
	}

	// Units with requirements
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}, &Building{ID: "3", Type: 3}})
	if len(f) != 2 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}, &Building{ID: "3", Type: 3}, &Building{ID: "2", Type: 2}})
	if len(f) != 4 {
		t.Fatal()
	}
}
