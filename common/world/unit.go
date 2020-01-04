// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

func (s *SetOfUnits) Len() int {
	return len(*s)
}

func (s *SetOfUnits) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfUnits) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfUnits) Add(u *Unit) {
	*s = append(*s, u)
	sort.Sort(s)
}

func (s *SetOfUnits) getFollowerIndex(o *Unit) int {
	for i, f := range *s {
		if f.Id == o.Id {
			return i
		}
	}
	return -1
}

func (s *SetOfUnits) removeFollower(index int) {
	lastIdx := len(*s) - 1
	if lastIdx != index {
		(*s)[index] = (*s)[index]
	}
	(*s)[lastIdx] = nil
	(*s) = (*s)[:lastIdx]
	sort.Sort(s)
}

func (s *SetOfUnits) Remove(u *Unit) {
	if idx := s.getFollowerIndex(u); idx >= 0 {
		s.removeFollower(idx)
	}
}

func (w *World) UnitGet(city, id uint64) *Unit {
	c := w.CityGet(city)
	if c != nil {
		return c.Unit(id)
	}
	return nil
}

func (s *SetOfUnitTypes) Len() int {
	return len(*s)
}

func (s *SetOfUnitTypes) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfUnitTypes) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfUnitTypes) Add(u *UnitType) {
	*s = append(*s, u)
	sort.Sort(s)
}

func (w *World) UnitGetType(id uint64) *UnitType {
	for _, ut := range w.Definitions.Units {
		if ut.Id == id {
			return ut
		}
	}
	return nil
}
