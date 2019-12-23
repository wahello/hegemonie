// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
)

func (w *World) CityGet(id uint64) *City {
	for _, c := range w.Cities {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func (w *World) CityCheck(id uint64) bool {
	return w.CityGet(id) != nil
}

func (w *World) CityCreate(id, loc uint64) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	c0 := w.CityGet(id)
	if c0 != nil {
		if c0.Deleted {
			c0.Deleted = false
			return nil
		} else {
			return errors.New("City already exists")
		}
	}

	c := City{Id: id, Cell: loc, Units: make([]uint64, 0)}
	w.Cities = append(w.Cities, c)
	return nil
}

func (w *World) CitySpawnUnit(idCity, idType uint64) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	c := w.CityGet(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := w.GetUnitType(idType)
	if t == nil {
		return errors.New("Unit type not found")
	}

	unit := Unit{Id: w.getNextId(), Health: t.Health, Type: t.Id, City: idCity, Cell: 0}
	w.Units = append(w.Units, unit)

	c.Units = append(c.Units, unit.Id)
	return nil
}

func (c *City) CityGetBuilding(id uint64) *Building {
	for _, b := range c.Buildings {
		if id == b.Id {
			return &b
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

	b := Building{Id: w.getNextId(), Type: idType}
	c.Buildings = append(c.Buildings, b)
	return nil
}

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

func (w *World) CityShow(userId, characterId, cityId uint64) (view CityView, err error) {
	w.rw.RLock()
	defer w.rw.RUnlock()

	pCity := w.CityGet(cityId)
	pChar := w.CharacterGet(characterId)

	if pCity == nil || pChar == nil {
		err = errors.New("Not Found")
	} else if pCity.Deputy != characterId && pCity.Owner != characterId {
		err = errors.New("Forbidden")
	} else if pChar.User != userId {
		err = errors.New("Forbidden")
	} else {
		view.Id = pCity.Id
		view.Name = pCity.Name
		view.Owner.Id = pCity.Owner
		view.Deputy.Id = pCity.Deputy
		view.Buildings = make([]BuildingView, 0, len(pCity.Buildings))
		view.Units = make([]UnitView, 0, len(pCity.Units))

		// Compute the modifiers
		for i := 0; i < ResourceMax; i++ {
			view.Production.Buildings.Mult[i] = 1.0
			view.Production.Knowledge.Mult[i] = 1.0
			view.Production.Troops.Mult[i] = 1.0
			view.Stock.Buildings.Mult[i] = 1.0
			view.Stock.Knowledge.Mult[i] = 1.0
			view.Stock.Troops.Mult[i] = 1.0
		}

		for _, b := range pCity.Buildings {
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
		for _, unitId := range pCity.Units {
			u := w.GetUnit(unitId)
			v := UnitView{}
			v.Id = u.Id
			v.Type = *w.GetUnitType(u.Type)
			view.Units = append(view.Units, v)
			for i := 0; i < ResourceMax; i++ {
				view.Production.Troops.Plus[i] += v.Type.Prod.Plus[i]
				view.Production.Troops.Mult[i] *= v.Type.Prod.Mult[i]
			}
		}

		// Apply all the modifiers on the production
		view.Production.Base = pCity.Production
		view.Production.Actual = pCity.Production
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
		view.Stock.Base = pCity.Stock
		view.Stock.Actual = pCity.Stock
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

		view.Stock.Usage = pCity.Stock
		return view, err
	}

	return view, err
}
