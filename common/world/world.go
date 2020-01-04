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
	w.Auth.Users = make(SetOfUsers, 0)
	w.Auth.Characters = make(SetOfCharacters, 0)
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

	if !sort.IsSorted(&w.Auth.Users) {
		return errors.New("user sequence: unsorted")
	}
	for i, u := range w.Auth.Users {
		if uint64(i)+1 != u.Id {
			return errors.New(fmt.Sprintf("user sequence: hole at %d", i))
		}
	}

	if !sort.IsSorted(&w.Auth.Characters) {
		return errors.New("character sequence: unsorted")
	}
	for i, c := range w.Auth.Characters {
		if uint64(i)+1 != c.Id {
			return errors.New(fmt.Sprintf("character sequence: hole at %d", i))
		}
	}

	if !sort.IsSorted(&w.Live.Cities) {
		return errors.New("city sequence: unsorted")
	}
	for i, c := range w.Live.Cities {
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

func (w *World) PostLoad() error {
	sort.Sort(&w.Auth.Users)
	sort.Sort(&w.Auth.Characters)
	sort.Sort(&w.Definitions.Knowledges)
	sort.Sort(&w.Definitions.Buildings)
	sort.Sort(&w.Definitions.Units)
	sort.Sort(&w.Live.Armies)
	sort.Sort(&w.Live.Cities)

	for _, a := range w.Live.Armies {
		sort.Sort(&a.Units)
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
