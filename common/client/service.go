// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hclient

import (
	"github.com/jfsmig/hegemonie/common/world"
)

type Region interface {
	Auth(args *AuthArgs, reply *AuthReply) error

	UserShow(args *UserShowArgs, reply *UserShowReply) error

	CharacterShow(args *CharacterShowArgs, reply *CharacterShowReply) error

	// Return an abstract of the City
	CityShow(args *CityShowArgs, reply *CityShowReply) error

	// Start the study of a new Knowledge
	CityStudy(args *CityStudyArgs, reply *CityStudyReply) error

	// Start the training of a new Unit
	CityTrain(args *CityTrainArgs, reply *CityTrainReply) error

	// Start a new building
	CityBuild(args *CityBuildArgs, reply *CityBuildReply) error

	// Create an empty army
	CityCreateArmy(args *CityCreateArmyArgs, reply *CityCreateArmyReply) error

	// Move in the given Army a Unit that must be ready and contained in the given City.
	CityTransferUnit(args *CityTransferUnitArgs, reply *CityTransferUnitReply) error

	// Push a new command for the Army.
	CityCommandArmy(args *CityCommandArmyArgs, reply *CityCommandArmyReply) error

	// Returns a string representation of the Map in the "dot" format
	// See https://graphviz.org for more information.
	MapDot(args *MapDotArgs, reply *MapDotReply) error

	// Perform an integrity check of the map
	MapCheck(args *MapCheckArgs, reply *MapCheckReply) error

	// Recompute secondary information when the definition of the Map has changed.
	MapRehash(args *MapRehashArgs, reply *MapRehashReply) error

	MapPlaces(args *MapPlacesArgs, reply *MapPlacesReply) error
	MapCities(args *MapCitiesArgs, reply *MapCitiesReply) error
	MapArmies(args *MapArmiesArgs, reply *MapArmiesReply) error

	// Persist the whole game information
	AdminSave(args *AdminSaveArgs, reply *AdminSaveReply) error

	// Perform an integrity check on the whole game information
	AdminCheck(args *AdminCheckArgs, reply *AdminCheckReply) error

	// Produce the resources in each City
	RoundProduce(args *RoundProduceArgs, reply *RoundProduceReply) error

	// Pley one movement step for each concerned Army
	RoundMove(args *RoundMoveArgs, reply *RoundMoveReply) error
}

type AuthArgs struct {
	UserMail string
	UserPass string
}

type AuthReply struct {
	Id uint64
}

type CityShowArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64
}

type CityShowReply struct {
	View world.CityView
}

type CityCommandArmyArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	ArmyId uint64
	Cell   uint64
	Action uint64
}

type CityCommandArmyReply struct{}

type CityBuildArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	BuildingId uint64
}

type CityBuildReply struct {
	Id uint64
}

type CityTrainArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	UnitId uint64
}

type CityTrainReply struct {
	Id uint64
}

type CityStudyArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	KnowledgeId uint64
}

type CityStudyReply struct {
	Id uint64
}

type CityCreateArmyArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	Name string
}

type CityCreateArmyReply struct {
	Id uint64
}

type CityTransferUnitArgs struct {
	UserId      uint64
	CharacterId uint64
	CityId      uint64

	UnitId uint64
	ArmyId uint64
}

type CityTransferUnitReply struct{}

type MapRehashArgs struct{}

type MapRehashReply struct{}

type MapCheckArgs struct{}

type MapCheckReply struct{}

type MapDotArgs struct{}

type MapDotReply struct {
	Dot string
}

type MapPlacesArgs struct{}

type MapPlacesReply struct {
	Items world.Map
}

type MapCitiesArgs struct{}

type MapCitiesReply struct {
	Items []*world.City
}

type MapArmiesArgs struct{}

type MapArmiesReply struct {
	Items []*world.Army
}

type AdminSaveArgs struct{}

type AdminSaveReply struct{}

type AdminCheckArgs struct{}

type AdminCheckReply struct{}

type RoundProduceArgs struct{}

type RoundProduceReply struct{}

type RoundMoveArgs struct{}

type RoundMoveReply struct{}

type UserShowArgs struct {
	UserId uint64
}

type UserShowReply struct {
	View world.UserView
}

type CharacterShowArgs struct {
	UserId      uint64
	CharacterId uint64
}

type CharacterShowReply struct {
	View world.CharacterView
}
