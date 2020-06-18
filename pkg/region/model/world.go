// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

func (w *World) WLock() { w.rw.Lock() }

func (w *World) WUnlock() { w.rw.Unlock() }

func (w *World) RLock() { w.rw.RLock() }

func (w *World) RUnlock() { w.rw.RUnlock() }

func (w *World) getNextId() uint64 {
	return atomic.AddUint64(&w.NextId, 1)
}

func (w *World) Produce() {
	for _, c := range w.Live.Cities {
		c.Produce(w)
	}
}

func (w *World) Move() {
	for _, c := range w.Live.Cities {
		for _, a := range c.Armies {
			a.Move(w)
		}
	}
}

func (w *World) UnitTypeGet(id uint64) *UnitType {
	return w.Definitions.Units.Get(id)
}

func (w *World) UnitGetFrontier(owned []*Building) []*UnitType {
	return w.Definitions.Units.Frontier(owned)
}

func (w *World) CityGet(id uint64) *City {
	return w.Live.Cities.Get(id)
}

func (w *World) CityCheck(id uint64) bool {
	return w.CityGet(id) != nil
}

func (w *World) CityCreateModel(loc uint64, model *City) (*City, error) {
	cell := w.Places.CellGet(loc)
	if cell == nil || cell.City != 0 {
		return nil, errors.New("Location already occupied")
	}

	id := w.getNextId()
	city := CopyCity(model)
	city.ID = id
	cell.City = id
	w.Live.Cities.Add(city)
	return city, nil
}

func (w *World) CityCreate(loc uint64) (*City, error) {
	return w.CityCreateModel(loc, nil)
}

func (w *World) CityCreateRandom(loc uint64) (*City, error) {
	if len(w.Config.CityPatterns) > 0 {
		i := rand.Intn(len(w.Config.CityPatterns))
		model := w.Config.CityPatterns[i]
		return w.CityCreateModel(loc, &model)
	} else {
		return w.CityCreateModel(loc, nil)
	}
}

func (w *World) CityGetAndCheck(characterId, cityId uint64) (*City, error) {
	// Fetch + sanity checks about the city
	pCity := w.CityGet(cityId)
	if pCity == nil {
		return nil, errors.New("Not Found")
	}
	if pCity.Deputy != characterId && pCity.Owner != characterId {
		return nil, errors.New("Forbidden")
	}

	return pCity, nil
}

func (w *World) Cities(idChar uint64) []*City {
	rep := make([]*City, 0)
	for _, c := range w.Live.Cities {
		if c.Owner == idChar || c.Deputy == idChar {
			rep = append(rep, c)
		}
	}
	return rep
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
