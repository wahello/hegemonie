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

func (w *World) UnitGet(id uint64) *Unit {
	for _, c := range w.Units {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (w *World) UnitGetType(id uint64) *UnitType {
	for _, c := range w.UnitTypes {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (u *Unit) insulate(w *World) {
	if u.Army != 0 {
		a := w.ArmyGet(u.Army)
		if a != nil {
			a.units.Remove(u)
		}
	}
	if u.City != 0 {
		c := w.CityGet(u.City)
		if c != nil {
			c.units.Remove(u)
		}
	}
	u.Army = 0
	u.City = 0
}

func (u *Unit) Incorporate(a *Army, w *World) {
	u.insulate(w)
	u.Army = a.Id
	a.units.Add(u)
}

func (u *Unit) Defend(c *City, w *World) {
	u.insulate(w)
	u.City = c.Id
	c.units.Add(u)
}
