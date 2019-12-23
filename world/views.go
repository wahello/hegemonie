// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

type NamedItem struct {
	Name string
	Id   uint64
}

// Read-Only hydrated version of a User
type UserView struct {
	Id         uint64
	Name       string
	Email      string
	Admin      bool
	Inactive   bool
	Characters []NamedItem
}

func (cv *UserView) GetNamed() NamedItem {
	return NamedItem{Id: cv.Id, Name: cv.Name}
}

// Read-Only hydrated version of a Character
type CharacterView struct {
	Id       uint64
	Name     string
	User     NamedItem
	OwnerOf  []NamedItem
	DeputyOf []NamedItem
}

func (cv *CharacterView) GetNamed() NamedItem {
	return NamedItem{Id: cv.Id, Name: cv.Name}
}

// Read-Only hydrated version of a Unit
type UnitView struct {
	Id   uint64
	Type UnitType
}

func (cv *UnitView) GetNamed() NamedItem {
	return NamedItem{Id: cv.Id, Name: cv.Type.Name}
}

// Read-Only hydrated version of a Building
type BuildingView struct {
	Id   uint64
	Type BuildingType
}

func (cv *BuildingView) GetNamed() NamedItem {
	return NamedItem{Id: cv.Id, Name: cv.Type.Name}
}

// Read-Only hydrated version of a City
type CityView struct {
	Id         uint64
	Name       string
	Cell       uint64
	Owner      NamedItem
	Deputy     NamedItem
	Production ProductionView
	Stock      StockView
	Units      []UnitView
	Buildings  []BuildingView
}

func (cv *CityView) GetNamed() NamedItem {
	return NamedItem{Id: cv.Id, Name: cv.Name}
}

type ProductionView struct {
	Base      Resources
	Knowledge ResourceModifiers
	Buildings ResourceModifiers
	Troops    ResourceModifiers
	Actual    Resources
}

type StockView struct {
	Base      Resources
	Knowledge ResourceModifiers
	Buildings ResourceModifiers
	Troops    ResourceModifiers
	Actual    Resources

	// Resources currently stored in the City, dispatched among all the building,
	// including the implicit town hall.
	Usage Resources
}