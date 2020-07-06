// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

func (w *World) WLock() { w.rw.Lock() }

func (w *World) WUnlock() { w.rw.Unlock() }

func (w *World) RLock() { w.rw.RLock() }

func (w *World) RUnlock() { w.rw.RUnlock() }

func (w *World) CreateRegion(name, mapName string) (*Region, error) {
	w.WLock()
	defer w.WUnlock()

	if w.Regions.Has(name) {
		return nil, errRegionExists
	}
	r := &Region{
		Name:    name,
		MapName: mapName,
		Cities:  make(SetOfCities, 0),
		Fights:  make(SetOfFights, 0),
		world:   w,
	}
	return r, nil
}

func (w *World) UnitTypeGet(id uint64) *UnitType {
	return w.Definitions.Units.Get(id)
}

func (w *World) UnitGetFrontier(owned []*Building) []*UnitType {
	return w.Definitions.Units.Frontier(owned)
}

func (w *World) BuildingTypeGet(id uint64) *BuildingType {
	return w.Definitions.Buildings.Get(id)
}

func (w *World) BuildingGetFrontier(pop int64, built []*Building, owned []*Knowledge) []*BuildingType {
	// TODO(jfs): Maybe speed the execution with a reverse index of Requires
	return w.Definitions.Buildings.Frontier(pop, built, owned)
}

func (w *World) KnowledgeTypeGet(id uint64) *KnowledgeType {
	return w.Definitions.Knowledges.Get(id)
}

func (w *World) KnowledgeGetFrontier(owned []*Knowledge) []*KnowledgeType {
	return w.Definitions.Knowledges.Frontier(owned)
}
