// Copyright (C) 2018-2020 Hegemonie's AUTHORS
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

	if !sort.IsSorted(&w.Auth.Characters) {
		return errors.New("character sequence: unsorted")
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
	// Sort all the lookup arrays
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
	for _, c := range w.Live.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
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
	if len(w.Auth.Users) > 0 {
		last := w.Auth.Users[len(w.Auth.Users)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
	if len(w.Auth.Characters) > 0 {
		last := w.Auth.Characters[len(w.Auth.Characters)-1]
		if last.Id > maxId {
			maxId = last.Id + 1
		}
	}
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
