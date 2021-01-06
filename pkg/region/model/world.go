// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

func (w *World) WLock() { w.rw.Lock() }

func (w *World) WUnlock() { w.rw.Unlock() }

func (w *World) RLock() { w.rw.RLock() }

func (w *World) RUnlock() { w.rw.RUnlock() }

func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

type NamedCity struct {
	Name string
	ID   uint64
}

func (w *World) CreateRegion(name, mapName string, cities []NamedCity) (*Region, error) {
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
	for _, x := range cities {
		city, err := r.CityCreate(x.ID)
		if err != nil {
			return nil, err
		}
		city.Name = x.Name
	}
	w.Regions.Add(r)
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
