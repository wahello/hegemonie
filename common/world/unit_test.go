// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sort"
	"testing"
)

func TestSetOfUnitType(t *testing.T) {
	u2 := &UnitType{Id: 2}
	s := SetOfUnitTypes{}
	s.Add(&UnitType{Id: 1})
	s.Add(&UnitType{Id: 3})
	s.Add(u2)
	s.Add(&UnitType{Id: 4})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get(2) {
		t.Fatal()
	}
}

func TestSetOfUnit(t *testing.T) {
	u2 := &Unit{Id: 2, Type: 1}
	s := SetOfUnits{}
	s.Add(&Unit{Id: 1, Type: 1})
	s.Add(&Unit{Id: 3, Type: 1})
	s.Add(u2)
	s.Add(&Unit{Id: 4, Type: 1})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get(2) {
		t.Fatal()
	}
}

func TestUnitFrontier(t *testing.T) {
	ut := SetOfUnitTypes{}
	ut.Add(&UnitType{Id: 1})
	ut.Add(&UnitType{Id: 2, RequiredBuilding: 2})
	ut.Add(&UnitType{Id: 3, RequiredBuilding: 2})
	ut.Add(&UnitType{Id: 4, RequiredBuilding: 3})

	b := SetOfBuildings{}
	b.Add(&Building{Id: 1, Type: 1})
	b.Add(&Building{Id: 2, Type: 2})
	b.Add(&Building{Id: 3, Type: 3})

	var f []*UnitType

	// Units without requirement
	f = ut.Frontier(SetOfBuildings{})
	if len(f) != 1 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{Id: 1, Type: 1}})
	if len(f) != 1 {
		t.Fatal()
	}

	// Units with requirements
	f = ut.Frontier(SetOfBuildings{&Building{Id: 1, Type: 1}, &Building{Id: 3, Type: 3}})
	if len(f) != 2 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{Id: 1, Type: 1}, &Building{Id: 3, Type: 3}, &Building{Id: 2, Type: 2}})
	if len(f) != 4 {
		t.Fatal()
	}
}
