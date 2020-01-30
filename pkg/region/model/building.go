// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import "sort"

func (s SetOfBuildings) Len() int           { return len(s) }
func (s SetOfBuildings) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfBuildings) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfBuildings) Get(id uint64) *Building {
	for _, b := range s {
		if b.Id == id {
			return b
		}
	}
	return nil
}

func (s *SetOfBuildings) Add(b *Building) {
	*s = append(*s, b)
	sort.Sort(s)
}

func (s SetOfBuildingTypes) Len() int           { return len(s) }
func (s SetOfBuildingTypes) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfBuildingTypes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfBuildingTypes) Get(id uint64) *BuildingType {
	for _, bt := range s {
		if bt.Id == id {
			return bt
		}
	}
	return nil
}

func (s *SetOfBuildingTypes) Add(b *BuildingType) {
	*s = append(*s, b)
	sort.Sort(s)
}

func (s SetOfBuildingTypes) Slice(marker uint64, max uint32) []*BuildingType {
	start := sort.Search(len(s), func(i int) bool { return s[i].Id > marker })
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfBuildingTypes) Frontier(pop int64, built []*Building, owned []*Knowledge) []*BuildingType {
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
		if bt.PopRequired > pop {
			return false
		}
		if bt.Unique && bmap[bt.Id] {
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
	for _, bt := range s {
		if valid(bt) {
			result = append(result, bt)
		}
	}
	return result
}

func (w *World) BuildingGetFrontier(pop int64, built []*Building, owned []*Knowledge) []*BuildingType {
	// TODO(jfs): Maybe speed the execution with a reverse index of Requires
	return w.Definitions.Buildings.Frontier(pop, built, owned)
}

func (w *World) BuildingTypeGet(id uint64) *BuildingType {
	return w.Definitions.Buildings.Get(id)
}
