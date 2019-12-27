// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
)

func (s *SetOfCities) Len() int {
	return len(*s)
}

func (s *SetOfCities) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfCities) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (w *World) CityGet(id uint64) *City {
	for _, c := range w.Cities {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (w *World) CityCheck(id uint64) bool {
	return w.CityGet(id) != nil
}

func (w *World) CityCreate(loc uint64) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	c := City{Id: w.getNextId(), Cell: loc, units: make(SetOfUnits, 0)}
	w.Cities = append(w.Cities, &c)
	return c.Id, nil
}

func (w *World) CitySpawnUnit(idCity, idType uint64) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	c := w.CityGet(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := w.UnitGetType(idType)
	if t == nil {
		return errors.New("Unit type not found")
	}

	unit := &Unit{Id: w.getNextId(), Health: t.Health, Type: t.Id, City: idCity}
	w.Units.Add(unit)
	unit.Defend(c, w)
	return nil
}

func (c *City) CityGetBuilding(id uint64) *Building {
	for _, b := range c.buildings {
		if id == b.Id {
			return b
		}
	}
	return nil
}

func (w *World) CitySpawnBuilding(idCity, idType uint64) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	c := w.CityGet(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := w.GetBuildingType(idType)
	if t == nil {
		return errors.New("Building tye not found")
	}

	// TODO(jfs): consume the resources

	b := &Building{Id: w.getNextId(), Type: idType}
	c.buildings.Add(b)
	return nil
}

func (w *World) CityShow(userId, characterId, cityId uint64) (view CityView, err error) {
	w.rw.RLock()
	defer w.rw.RUnlock()

	// Fetch + sanity checks about the city
	pCity := w.CityGet(cityId)
	if pCity == nil {
		err = errors.New("Not Found")
		return
	}
	if pCity.Deputy != characterId && pCity.Owner != characterId {
		err = errors.New("Forbidden")
		return
	}

	// Fetch + senity checks about the City
	pOwner := w.CharacterGet(pCity.Owner)
	pDeputy := w.CharacterGet(pCity.Deputy)
	if pOwner == nil || pDeputy == nil {
		err = errors.New("Not Found")
		return
	}
	if pOwner.User != userId && pDeputy.User != userId {
		err = errors.New("Forbidden")
		return
	}

	view = pCity.Show(w)
	view.Owner.Name = pOwner.Name
	view.Deputy.Name = pDeputy.Name
	return
}

func (c *City) Show(w *World) (view CityView) {
	view.Id = c.Id
	view.Name = c.Name
	view.Owner.Id = c.Owner
	view.Deputy.Id = c.Deputy
	view.Buildings = make([]BuildingView, 0, len(c.buildings))
	view.Units = make([]UnitView, 0, len(c.units))
	view.Armies = make([]NamedItem, 0)

	for _, a := range c.armies {
		view.Armies = append(view.Armies, NamedItem{Id: a.Id, Name: a.Name})
	}

	// Compute the modifiers
	for i := 0; i < ResourceMax; i++ {
		view.Production.Buildings.Mult[i] = 1.0
		view.Production.Knowledge.Mult[i] = 1.0
		view.Production.Troops.Mult[i] = 1.0
		view.Stock.Buildings.Mult[i] = 1.0
		view.Stock.Knowledge.Mult[i] = 1.0
		view.Stock.Troops.Mult[i] = 1.0
	}

	for _, b := range c.buildings {
		v := BuildingView{}
		v.Id = b.Id
		v.Type = *w.GetBuildingType(b.Type)
		view.Buildings = append(view.Buildings, v)
		for i := 0; i < ResourceMax; i++ {
			view.Production.Buildings.Plus[i] += v.Type.Prod.Plus[i]
			view.Production.Buildings.Mult[i] *= v.Type.Prod.Mult[i]
			view.Stock.Buildings.Plus[i] += v.Type.Stock.Plus[i]
			view.Stock.Buildings.Mult[i] *= v.Type.Stock.Mult[i]
		}
	}
	for _, u := range c.units {
		v := UnitView{}
		v.Id = u.Id
		v.Type = *w.UnitGetType(u.Type)
		view.Units = append(view.Units, v)
		for i := 0; i < ResourceMax; i++ {
			view.Production.Troops.Plus[i] += v.Type.Prod.Plus[i]
			view.Production.Troops.Mult[i] *= v.Type.Prod.Mult[i]
		}
	}

	// Apply all the modifiers on the production
	view.Production.Base = c.Production
	view.Production.Actual = c.Production
	for i := 0; i < ResourceMax; i++ {
		v := float64(view.Production.Base[i])
		v = v * view.Production.Troops.Mult[i]
		v = v * view.Production.Buildings.Mult[i]
		v = v * view.Production.Knowledge.Mult[i]

		vi := int64(v)
		vi = vi + view.Production.Troops.Plus[i]
		vi = vi + view.Production.Buildings.Plus[i]
		vi = vi + view.Production.Knowledge.Plus[i]

		view.Production.Actual[i] = uint64(vi)
	}

	// Apply all the modifiers on the stock
	view.Stock.Base = c.StockCapacity
	view.Stock.Actual = c.StockCapacity
	view.Stock.Usage = c.Stock
	for i := 0; i < ResourceMax; i++ {
		v := float64(view.Stock.Base[i])
		v = v * view.Stock.Troops.Mult[i]
		v = v * view.Stock.Buildings.Mult[i]
		v = v * view.Stock.Knowledge.Mult[i]

		vi := int64(v)
		vi = vi + view.Stock.Troops.Plus[i]
		vi = vi + view.Stock.Buildings.Plus[i]
		vi = vi + view.Stock.Knowledge.Plus[i]

		view.Stock.Actual[i] = uint64(vi)
	}

	return
}

func (c *City) Produce(w *World) {
	// Pre-compute the modified values of Stock and Production
	view := c.Show(w)
	post := view.Stock.Usage

	post.Add(&view.Production.Actual)

	for _, b := range c.buildings {
		if b.Ticks > 0 {
			bt := w.GetBuildingType(b.Id)
			if post.GreaterOrEqualTo(&bt.Cost) {
				post.Remove(&bt.Cost)
				b.Ticks--
			}
		}
	}

	for _, u := range c.units {
		if u.Ticks > 0 {
			ut := w.UnitGetType(u.Type)
			if post.GreaterOrEqualTo(&ut.Cost) {
				post.Remove(&ut.Cost)
				u.Ticks--
			}
		}
	}

	post.TrimTo(&view.Stock.Actual)
	c.Stock = post
}
