// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync/atomic"
)

func (w *World) Init() {
	w.rw.Lock()
	defer w.rw.Unlock()

	w.Places.Init()

	if w.NextId <= 0 {
		w.NextId = 1
	}
	w.Live.Armies = make(SetOfArmies, 0)
	w.Live.Cities = make(SetOfCities, 0)
	w.Definitions.Units = make(SetOfUnitTypes, 0)
	w.Definitions.Buildings = make(SetOfBuildingTypes, 0)
	w.Definitions.Knowledges = make(SetOfKnowledgeTypes, 0)
}

func (w *World) Check() error {
	var err error
	err = w.Places.Check(w)
	if err != nil {
		return err
	}

	if !sort.IsSorted(&w.Definitions.Knowledges) {
		return errors.New("knowledge sequence: unsorted")
	}

	if !sort.IsSorted(&w.Definitions.Buildings) {
		return errors.New("building sequence: unsorted")
	}

	if !sort.IsSorted(&w.Definitions.Units) {
		return errors.New("unit sequence: unsorted")
	}

	if !sort.IsSorted(&w.Live.Cities) {
		return errors.New("city sequence: unsorted")
	}

	if !sort.IsSorted(&w.Live.Armies) {
		return errors.New("army sequence: unsorted")
	}

	for _, a := range w.Live.Armies {
		if !sort.IsSorted(&a.Units) {
			return errors.New("unit sequence: unsorted")
		}
	}
	for _, a := range w.Live.Cities {
		if !sort.IsSorted(&a.Knowledges) {
			return errors.New("knowledge sequence: unsorted")
		}
		if !sort.IsSorted(&a.Buildings) {
			return errors.New("building sequence: unsorted")
		}
		if !sort.IsSorted(&a.Units) {
			return errors.New("unit sequence: unsorted")
		}
	}

	return nil
}

func (w *World) WLock() { w.rw.Lock() }

func (w *World) WUnlock() { w.rw.Unlock() }

func (w *World) RLock() { w.rw.RLock() }

func (w *World) RUnlock() { w.rw.RLock() }

func (w *World) getNextId() uint64 {
	return atomic.AddUint64(&w.NextId, 1)
}

func (w *World) DumpJSON(dst io.Writer) error {
	return json.NewEncoder(dst).Encode(w)
}

func (w *World) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&w.Definitions.Knowledges)
	sort.Sort(&w.Definitions.Buildings)
	sort.Sort(&w.Definitions.Units)
	sort.Sort(&w.Live.Armies)
	sort.Sort(&w.Live.Cities)
	for _, a := range w.Live.Armies {
		sort.Sort(&a.Units)
	}
	for _, c := range w.Live.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
	}

	if err := w.Live.Armies.Check(); err != nil {
		return err
	}
	if err := w.Live.Cities.Check(); err != nil {
		return err
	}

	// Link Armies and Cities
	for _, a := range w.Live.Armies {
		if a.City == 0 {
			return errors.New(fmt.Sprintf("Army %v points to no City", a))
		} else if c := w.CityGet(a.City); c == nil {
			return errors.New(fmt.Sprintf("Army %v points to ghost City", a))
		} else {
			c.armies.Add(a)
		}
	}

	// Compute the highest unique ID
	maxId := w.NextId
	if len(w.Definitions.Units) > 0 {
		last := w.Definitions.Units[len(w.Definitions.Units)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	if len(w.Definitions.Buildings) > 0 {
		last := w.Definitions.Buildings[len(w.Definitions.Buildings)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	if len(w.Definitions.Knowledges) > 0 {
		last := w.Definitions.Knowledges[len(w.Definitions.Knowledges)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	if len(w.Live.Armies) > 0 {
		last := w.Live.Armies[len(w.Live.Armies)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	if len(w.Live.Cities) > 0 {
		last := w.Live.Cities[len(w.Live.Cities)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	for _, c := range w.Live.Cities {
		if len(c.Knowledges) > 0 {
			last := c.Knowledges[len(c.Knowledges)-1]
			if last.Id > maxId {
				maxId = last.Id + 1
			}
		}
		if len(c.Units) > 0 {
			last := c.Units[len(c.Units)-1]
			if last.Id > maxId {
				maxId = last.Id + 1
			}
		}
		if len(c.Buildings) > 0 {
			last := c.Buildings[len(c.Buildings)-1]
			if last.Id > maxId {
				maxId = last.Id + 1
			}
		}
	}
	for _, a := range w.Live.Armies {
		if len(a.Units) > 0 {
			last := a.Units[len(a.Units)-1]
			if last.Id > maxId {
				maxId = last.Id + 1
			}
		}
	}

	w.NextId = maxId
	return nil
}

func (w *World) Produce() {
	w.rw.Lock()
	defer w.rw.Unlock()

	for _, c := range w.Live.Cities {
		c.Produce(w)
	}
}

func (w *World) Move() {
	w.rw.Lock()
	defer w.rw.Unlock()

	for _, a := range w.Live.Armies {
		a.Move(w)
	}
}

func (w *World) UnitTypeGet(id uint64) *UnitType {
	return w.Definitions.Units.Get(id)
}

func (w *World) UnitGetFrontier(owned []*Building) []*UnitType {
	return w.Definitions.Units.Frontier(owned)
}

func (w *World) ArmyCreate(c *City, name string) (*Army, error) {
	a := &Army{
		Id: w.getNextId(), City: c.Id, Cell: c.Cell,
		Name: name, Units: make(SetOfUnits, 0),
		Targets: make([]Command, 0),
	}
	w.Live.Armies.Add(a)
	c.armies.Add(a)
	return a, nil
}

func (w *World) ArmyGet(id uint64) *Army {
	return w.Live.Armies.Get(id)
}

func (w *World) CityGet(id uint64) *City {
	return w.Live.Cities.Get(id)
}

func (w *World) CityCheck(id uint64) bool {
	return w.CityGet(id) != nil
}

func (w *World) CityCreate(loc uint64) (uint64, error) {
	id := w.getNextId()
	w.Live.Cities.Create(id, loc)
	return id, nil
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
	return rep[:]
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
