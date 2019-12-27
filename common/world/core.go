// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sync"
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

type Army struct {
	// The unique ID of the current Army
	Id uint64

	// The ID of the City that controls the current Army
	City uint64

	// The ID of the Cell the Army is on
	Cell uint64

	// The goal
	Target uint64

	// A display name for the current City
	Name string

	units SetOfUnits
}

type BuildingType struct {
	// Unique ID of the BuildingType
	Id uint64

	// Display name of the current BuildingType
	Name string

	// How many ticks for the construction
	Ticks uint

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

	// How many construction rounds remain before the building's achievement
	Ticks uint
}

type Character struct {
	// The unique identifier of the current Character
	Id uint64

	// The unique identifier of the only User that controls the Character.
	User uint64

	// The display name of the current Character
	Name string
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

	// PRIVATE
	buildings SetOfBuildings
	units     SetOfUnits
	armies    SetOfArmies
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

	// How many ticks
	Ticks uint

	// Instantiation cost of the current UnitType
	Cost Resources

	// Might positive (resource boost) or more commonly negative (maintenance cost)
	Prod ResourceModifiers
}

// Both Cell and City must not be 0, and have a non-0 value
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint

	// A copy of the definition for the current Unit.
	Type uint64

	// The unique Id of the map cell the current Unit is on.
	Army uint64

	// The unique Id of the City the Unit is in.
	City uint64

	// How many ticks remain before the Troop training is finished
	Ticks uint
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

type SetOfArmies []*Army

type SetOfUnits []*Unit

type SetOfUsers []*User

type SetOfBuildings []*Building

type SetOfCharacters []*Character

type SetOfCities []*City

type World struct {
	Armies        SetOfArmies
	Users         SetOfUsers
	Characters    SetOfCharacters
	Cities        SetOfCities
	Units         SetOfUnits
	UnitTypes     []*UnitType
	BuildingTypes []*BuildingType

	NextId uint64
	Salt   string
	rw     sync.RWMutex

	Places Map
}
