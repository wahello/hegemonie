// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

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

func (w *World) GetBuildingType(id uint64) *BuildingType {
	for _, i := range w.BuildingTypes {
		if i.Id == id {
			return i
		}
	}
	return nil
}
