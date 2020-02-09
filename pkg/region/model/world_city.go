// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
)

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
