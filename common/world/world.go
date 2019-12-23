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

const (
	ResourceMax = 4
)

type Resources [ResourceMax]uint64

type ResourcesIncrement [ResourceMax]int64

type ResourcesMultiplier [ResourceMax]float64

type ResourceModifiers struct {
	Mult ResourcesMultiplier
	Plus ResourcesIncrement
}

type UnitType struct {
	// Unique Id of the Unit Type
	Id uint64

	// The number of health point for that type of unit.
	// A health equal to 0 means the death of the unit.
	Health uint

	// How affected is that type of unit by a loss of Health.
	// Must be between 0 and 1.
	// 0 means that the capacity of the Unit isn't affected by a health reduction.
	// 1 means that the capacity of the Unit loses an equal percentage of its capacity
	// for a loss of health (in other words, a HealthFactor of 1 means that the Unit
	// will hit at 90% of its maximal power if it has 90% of its health points).
	HealthFactor float64

	// The display name of the Unit Type
	Name string

	// Instantiation cost of the current UnitType
	Cost Resources

	// Might positive (resource boost) or more commonly negative (maintenance cost)
	Prod ResourceModifiers
}

// Both Cell and City must not be 0, and have a non-0 value.
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint

	// A copy of the definition for the current Unit.
	Type uint64

	// The unique Id of the map cell the current Unit is on.
	Cell uint64

	// The unique Id of the City the Unit is in.
	City uint64
}

type BuildingType struct {
	// Unique ID of the BuildingType
	Id uint64

	// Display name of the current BuildingType
	Name string

	// How much does the production cost
	Cost Resources

	// Impat of the current Building on the total storage capacity of the City.
	Stock ResourceModifiers

	// Increment of resources produced by this building.
	Prod ResourceModifiers
}

type Building struct {
	// The unique ID of the current Building
	Id uint64

	// The unique ID of the BuildingType associated to the current Building
	Type uint64
}

type City struct {
	// The unique ID of the current City
	Id uint64

	// The unique ID of the main Character in charge of the City.
	// The Manager may name a Deputy manager in the City.
	Owner uint64

	// The unique ID of a second Character in charge of the City.
	Deputy uint64

	// The unique ID of the Cell the current City is built on.
	// This is redundant with the City field in the Cell structure.
	// Both information must match.
	Cell uint64

	// The display name of the current City
	Name string

	// Resources stock owned by the current City
	Stock Resources

	// Maximum amounts of each resources that might be stored in the town hall
	// of the city. That limit doesn't consider the modifiers.
	StockCapacity Resources

	// Resources produced each round by the City, before the enforcing of
	// Production Boosts ans Production Multipliers
	Production Resources

	// Is the city still usable
	Deleted bool

	// An array of Units guarding the current City.
	// This is redundant with the City field of the Unit type.
	// Consider it as an index.
	Units []uint64

	Buildings []Building
}

type Character struct {
	// The unique identifier of the current Character
	Id uint64

	// The unique identifier of the only User that controls the Character.
	User uint64

	// The display name of the current Character
	Name string
}

type User struct {
	// The unique identifier of the current User
	Id uint64

	// The display name of the current User
	Name string

	// The unique email that authenticates the User.
	Email string

	// The hashed password that authenticates the User
	Password string

	// Has the current User the permission to manage the service.
	Admin bool `json:",omitempty"`

	// Can the user still login.
	Inactive bool `json:",omitempty"`
}

type SetOfUsers []User

type SetOfCharacters []Character

type SetOfCities []City

type World struct {
	Users         SetOfUsers
	Characters    SetOfCharacters
	Cities        SetOfCities
	Units         []Unit
	UnitTypes     []UnitType
	BuildingTypes []BuildingType

	NextId uint64
	Salt   string
	rw     sync.RWMutex

	Places Map
}

func (w *World) Init() {
	w.rw.Lock()
	defer w.rw.Unlock()

	w.Places.Init()

	if w.NextId <= 0 {
		w.NextId = 1
	}
	w.Users = make([]User, 0)
	w.Characters = make([]Character, 0)
	w.Cities = make([]City, 0)
	w.Units = make([]Unit, 0)
	w.UnitTypes = make([]UnitType, 0)
	w.BuildingTypes = make([]BuildingType, 0)
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
	sort.Sort(&w.Users)
	sort.Sort(&w.Characters)
	sort.Sort(&w.Cities)
	return nil
}
