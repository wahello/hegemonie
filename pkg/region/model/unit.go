// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

func (w *World) UnitGet(city, id uint64) *Unit {
	c := w.CityGet(city)
	if c != nil {
		return c.Unit(id)
	}
	return nil
}

func (s SetOfUnitTypes) Frontier(owned []*Building) []*UnitType {
	bIndex := make(map[uint64]bool)
	for _, b := range owned {
		bIndex[b.Type] = true
	}
	result := make([]*UnitType, 0)
	for _, ut := range s {
		if ut.RequiredBuilding == 0 || bIndex[ut.RequiredBuilding] {
			result = append(result, ut)
		}
	}
	return result
}

func (w *World) UnitTypeGet(id uint64) *UnitType {
	return w.Definitions.Units.Get(id)
}

func (w *World) UnitGetFrontier(owned []*Building) []*UnitType {
	return w.Definitions.Units.Frontier(owned)
}
