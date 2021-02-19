// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import "github.com/juju/errors"

// WLock acquires an exclusive ("writer") lock on the current world
// The request for a writer's lock has the priority on any equest for a reader's lock.
// It also forbid any concurrent writer or reader lock to be acquired.
func (w *World) WLock() { w.rw.Lock() }

// WUnlock releases an exclusive ("writer") lock on the current world
func (w *World) WUnlock() { w.rw.Unlock() }

// RLock acquires a shared("reader") lock on the current world
// Several requests for a reader's lock may be granted simultaneously but they are
// incompatible with any writer's locks. Their usage is reserved for read-nonly operations.
func (w *World) RLock() { w.rw.RLock() }

// RUnlock releases a shared("reader") lock on the current world
func (w *World) RUnlock() { w.rw.RUnlock() }

// SetNotifier changes the event notifier associated with the current World.
// That Notifier is user for inègame notifications, to collect a per-Character log.
func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

func (w *World) SetMapClient(c MapView) {
	// TODO(jfs): maybe wrap with a cache if perf is too poor
	w.mapView = c
}

// NamedCity is the information that a City with named Name should exist at the position ID on the graph map.
type NamedCity struct {
	Name string
	ID   uint64
}

func (w *World) CreateRegion(name, mapName string, cities []NamedCity) (*Region, error) {
	if w.Regions.Has(name) {
		return nil, errors.AlreadyExistsf("region found with id [%s]", name)
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
