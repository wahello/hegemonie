// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

func (s *SetOfBuildingTypes) Len() int {
	return len(*s)
}

func (s *SetOfBuildingTypes) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfBuildingTypes) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfBuildingTypes) Add(b *BuildingType) {
	*s = append(*s, b)
	sort.Sort(s)
}

func (w *World) BuildingTypeGet(id uint64) *BuildingType {
	for _, i := range w.Definitions.Buildings {
		if i.Id == id {
			return i
		}
	}
	return nil
}

func (s *SetOfBuildings) Len() int {
	return len(*s)
}

func (s *SetOfBuildings) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfBuildings) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfBuildings) Add(b *Building) {
	*s = append(*s, b)
	sort.Sort(s)
}

// TODO(jfs): Maybe speed the execution with a reverse index of Requires
func (w *World) BuildingGetFrontier(built []*Building, owned []*Knowledge) []*BuildingType {
	bmap := make(map[uint64]bool)
	pending := make(map[uint64]bool)
	finished := make(map[uint64]bool)
	for _, k := range owned {
		if k.Ticks == 0 {
			finished[k.Type] = true
		} else {
			pending[k.Type] = true
		}
	}
	for _, b := range built {
		bmap[b.Type] = true
	}

	valid := func(bt *BuildingType) bool {
		// TODO(jfs): Manage not only singleton
		if bmap[bt.Id] {
			return false
		}
		for _, c := range bt.Conflicts {
			if finished[c] || pending[c] {
				return false
			}
		}
		for _, c := range bt.Requires {
			if !finished[c] {
				return false
			}
		}
		return true
	}

	result := make([]*BuildingType, 0)
	for _, bt := range w.Definitions.Buildings {
		if valid(bt) {
			result = append(result, bt)
		}
	}
	return result
}
