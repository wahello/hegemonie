// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hege_world_client

import (
	"net/rpc"
)

type RegionClientTcp struct {
	endpoint string
}

const (
	kPrefix = "Region."
)

func DialClientTcp(endpoint string) *RegionClientTcp {
	srv := RegionClientTcp{}
	srv.endpoint = endpoint
	return &srv
}

func (s *RegionClientTcp) Auth(args *AuthArgs, reply *AuthReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"Auth", args, reply)
}

func (s *RegionClientTcp) UserShow(args *UserShowArgs, reply *UserShowReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"UserShow", args, reply)
}

func (s *RegionClientTcp) CharacterShow(args *CharacterShowArgs, reply *CharacterShowReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"CharacterShow", args, reply)
}

func (s *RegionClientTcp) CityShow(args *CityShowArgs, reply *CityShowReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"CityShow", args, reply)
}

func (s *RegionClientTcp) CityStudy(args *CityStudyArgs, reply *CityStudyReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"CityStudy", args, reply)
}

func (s *RegionClientTcp) CityTrain(args *CityTrainArgs, reply *CityTrainReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"CityTrain", args, reply)
}

func (s *RegionClientTcp) CityBuild(args *CityBuildArgs, reply *CityBuildReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"CityBuild", args, reply)
}

func (s *RegionClientTcp) MapDot(args *MapDotArgs, reply *MapDotReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapDot", args, reply)
}

func (s *RegionClientTcp) MapCheck(args *MapCheckArgs, reply *MapCheckReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapCheck", args, reply)
}

func (s *RegionClientTcp) MapRehash(args *MapRehashArgs, reply *MapRehashReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapRehash", args, reply)
}

func (s *RegionClientTcp) MapPlaces(args *MapPlacesArgs, reply *MapPlacesReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapPlaces", args, reply)
}

func (s *RegionClientTcp) MapCities(args *MapCitiesArgs, reply *MapCitiesReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapCities", args, reply)
}

func (s *RegionClientTcp) MapArmies(args *MapArmiesArgs, reply *MapArmiesReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"MapArmies", args, reply)
}

func (s *RegionClientTcp) AdminSave(args *AdminSaveArgs, reply *AdminSaveReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"AdminSave", args, reply)
}

func (s *RegionClientTcp) AdminCheck(args *AdminCheckArgs, reply *AdminCheckReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"AdminCheck", args, reply)
}

func (s *RegionClientTcp) RoundProduce(args *RoundProduceArgs, reply *RoundProduceReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"RoundProduce", args, reply)
}

func (s *RegionClientTcp) RoundMove(args *RoundMoveArgs, reply *RoundMoveReply) error {
	cnx, err := rpc.DialHTTP("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer cnx.Close()
	return cnx.Call(kPrefix+"RoundMove", args, reply)
}
