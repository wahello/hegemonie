// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

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
		if !bt.MultipleAllowed && bmap[bt.ID] {
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
