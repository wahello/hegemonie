// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

func (w *World) Init() {
	w.rw.Lock()
	defer w.rw.Unlock()

	w.Places.Init()

	if w.NextId <= 0 {
		w.NextId = 1
	}
	w.Armies = make([]*Army, 0)
	w.Users = make([]*User, 0)
	w.Characters = make([]*Character, 0)
	w.Cities = make([]*City, 0)
	w.Units = make([]*Unit, 0)
	w.UnitTypes = make([]*UnitType, 0)
	w.BuildingTypes = make([]*BuildingType, 0)
}

func (w *World) Check() error {
	var err error
	err = w.Places.Check(w)
	if err != nil {
		return err
	}

	if !sort.IsSorted(&w.Users) {
		return errors.New("user sequence: unsorted")
	}
	for i, u := range w.Users {
		if uint64(i)+1 != u.Id {
			return errors.New(fmt.Sprintf("user sequence: hole at %d", i))
		}
	}

	if !sort.IsSorted(&w.Characters) {
		return errors.New("character sequence: unsorted")
	}
	for i, c := range w.Characters {
		if uint64(i)+1 != c.Id {
			return errors.New(fmt.Sprintf("character sequence: hole at %d", i))
		}
	}

	if !sort.IsSorted(&w.Cities) {
		return errors.New("city sequence: unsorted")
	}
	for i, c := range w.Cities {
		if uint64(i)+1 != c.Id {
			return errors.New(fmt.Sprintf("City sequence: hole at %d", i))
		}
	}

	return nil
}

func (w *World) ReadLocker() sync.Locker {
	return w.rw.RLocker()
}

func (w *World) getNextId() uint64 {
	return atomic.AddUint64(&w.NextId, 1)
}

func (w *World) DumpJSON(dst io.Writer) error {
	return json.NewEncoder(dst).Encode(w)
}

func (w *World) LoadJSON(src io.Reader) error {
	err := json.NewDecoder(src).Decode(w)
	if err != nil {
		return err
	}
	sort.Sort(&w.Armies)
	sort.Sort(&w.Users)
	sort.Sort(&w.Characters)
	sort.Sort(&w.Cities)
	sort.Sort(&w.Units)

	// Link Units and Armies
	for _, u := range w.Units {
		if u.Army != 0 {
			a := w.ArmyGet(u.Army)
			if a == nil {
				return errors.New(fmt.Sprintf("Unit %v points to ghost Army", u))
			} else {
				u.Incorporate(a, w)
			}
		} else if u.City != 0 {
			c := w.CityGet(u.City)
			if c == nil {
				return errors.New(fmt.Sprintf("Unit %v points to ghost City", u))
			} else {
				u.Defend(c, w)
			}
		} else {
			return errors.New(fmt.Sprintf("Unit %v points to no City and no Army", u))
		}
	}

	// Link Armies and Cities
	for _, a := range w.Armies {
		if a.City == 0 {
			return errors.New(fmt.Sprintf("Army %v points to no City", a))
		} else if c := w.CityGet(a.City); c == nil {
			return errors.New(fmt.Sprintf("Army %v points to ghost City", a))
		} else {
			c.armies.Add(a)
		}
	}

	return nil
}

func (w *World) Produce() {
	w.rw.Lock()
	defer w.rw.Unlock()

	for _, c := range w.Cities {
		c.Produce(w)
	}
}

func (w *World) Move() {
	w.rw.Lock()
	defer w.rw.Unlock()

	for _, a := range w.Armies {
		a.Move(w)
	}
}
