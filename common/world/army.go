// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

func (s *SetOfArmies) Len() int {
	return len(*s)
}

func (s *SetOfArmies) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfArmies) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfArmies) Add(a *Army) {
	*s = append(*s, a)
	sort.Sort(s)
}

func (w *World) ArmyGet(id uint64) *Army {
	for _, a := range w.Armies {
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (a *Army) Move(w *World) error {
	src := a.Cell
	dst := a.Target
	nxt, err := w.Places.NextStep(src, dst)
	if err != nil {
		return err
	}

	a.Cell = nxt
	if nxt == dst {
		// TODO(jfs) There is something to do
	}
	return nil
}
