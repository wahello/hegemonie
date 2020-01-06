// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sync"
)

const (
	ResourceMax = 6
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

	// A display name for the current City
	Name string

	// The ID of the City that controls the current Army
	City uint64 `json:",omitempty"`

	// The ID of the Cell the Army is on
	Cell uint64 `json:",omitempty"`

	// The IS of a Cell of the Map that is a goal of the current movement of the Army
	Target uint64 `json:",omitempty"`

	// How many resources are carried by that Army
	Stock Resources

	Units SetOfUnits
}

type KnowledgeType struct {
	Id        uint64
	Name      string
	Ticks     uint `json:",omitempty"`
	Cost      Resources
	Requires  []uint64
	Conflicts []uint64
}

type Knowledge struct {
	Id    uint64
	Type  uint64
	Ticks uint `json:",omitempty"`
}

type BuildingType struct {
	// Unique ID of the BuildingType
	Id uint64

	// Display name of the current BuildingType
	Name string

	// How many ticks for the construction
	Ticks uint `json:",omitempty"`

	// How much does the production cost
	Cost Resources

	// Has the building to be unique a the City
	Unique bool `json:",omitempty"`

	// Impat of the current Building on the total storage capacity of the City.
	Stock ResourceModifiers

	// Increment of resources produced by this building.
	Prod ResourceModifiers

	// A set of KnowledgeType ID that must all be present in a City to let that City start
	// this kind of building.
	Requires []uint64

	// A set of KnowledgeType ID that must all be absent in a City to let that City start
	// this kind of building.
	Conflicts []uint64
}

type Building struct {
	// The unique ID of the current Building
	Id uint64

	// The unique ID of the BuildingType associated to the current Building
	Type uint64

	// How many construction rounds remain before the building's achievement
	Ticks uint `json:",omitempty"`
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
	Deputy uint64 `json:",omitempty"`

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
	Deleted bool `json:",omitempty"`

	// Tells if the City is in automatic mode.
	// The "auto" mode is intented for inactive or absent players.
	// The armies come home to defend the City, no new building or unit is spawned.
	// In the plans: a conservative behavior should be automated
	Auto bool `json:",omitempty"`

	Knowledges SetOfKnowledges

	Buildings SetOfBuildings

	// Units directly defending the current City
	Units SetOfUnits

	// PRIVATE
	// Armies under the responsibility of the current City
	armies SetOfArmies
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

	// A UnitType is only dependant on the presence of a Building of that BuildingType.
	RequiredBuilding uint64
}

// Both Cell and City must not be 0, and have a non-0 value
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// A copy of the definition for the current Unit.
	Type uint64

	// How many ticks remain before the Troop training is finished
	Ticks uint

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint `json:"H,omitempty"`
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

// A MapVertex is a Directed Vertex in the transport graph
type MapVertex struct {
	// Unique identifier of the source Cell
	S uint64

	// Unique identifier of the destination Cell
	D uint64

	// May the road be used by Units
	Deleted bool `json:",omitempty"`
}

// A MapNode is a Node is the directed graph of the transport network.
type MapNode struct {
	// The unique identifier of the current cell.
	Id uint64

	// Biome in which the cell is
	Biome uint64

	// The unique ID of the city present at this location.
	City uint64 `json:",omitempty"`
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	Cells  SetOfNodes
	Roads  SetOfVertices
	NextId uint64

	steps      map[vector]uint64
	dirtyRoads bool
	dirtyCells bool
	rw         sync.RWMutex
}

type SetOfArmies []*Army

type SetOfUnits []*Unit

type SetOfUnitTypes []*UnitType

type SetOfUsers []*User

type SetOfBuildings []*Building

type SetOfBuildingTypes []*BuildingType

type SetOfKnowledges []*Knowledge

type SetOfKnowledgeTypes []*KnowledgeType

type SetOfCharacters []*Character

type SetOfCities []*City

type SetOfNodes []MapNode

type SetOfVertices []MapVertex

type AuthBase struct {
	Users      SetOfUsers
	Characters SetOfCharacters
}

type DefinitionsBase struct {
	Units      SetOfUnitTypes
	Buildings  SetOfBuildingTypes
	Knowledges SetOfKnowledgeTypes
}

type LiveBase struct {
	Armies SetOfArmies
	Cities SetOfCities
}

type World struct {
	Auth        AuthBase
	Definitions DefinitionsBase
	Live        LiveBase

	NextId uint64
	Salt   string
	rw     sync.RWMutex

	Places Map
}
