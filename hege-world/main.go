// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	. "github.com/jfsmig/hegemonie/common/client"
	. "github.com/jfsmig/hegemonie/common/world"
)

var (
	pathSave string
)

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}

func save(w *World) error {
	if pathSave == "" {
		return errors.New("No save path configured")
	}
	p := pathSave + "/" + makeSaveFilename()
	p = filepath.Clean(p)
	out, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = w.DumpJSON(out)
	out.Close()
	if err != nil {
		_ = os.Remove(p)
		return err
	}

	latest := pathSave + "/latest"
	_ = os.Remove(latest)
	_ = os.Symlink(p, latest)
	return nil
}

type RegionService struct {
	w *World
}

func (s *RegionService) Auth(args *AuthArgs, reply *AuthReply) error {
	id, err := s.w.UserAuth(args.UserMail, args.UserPass)
	reply.Id = id
	return err
}

func (s *RegionService) UserShow(args *UserShowArgs, reply *UserShowReply) error {
	v, err := s.w.UserShow(args.UserId)
	reply.View = v
	return err
}

func (s *RegionService) CharacterShow(args *CharacterShowArgs, reply *CharacterShowReply) error {
	v, err := s.w.CharacterShow(args.UserId, args.CharacterId)
	reply.View = v
	return err
}

func (s *RegionService) CityShow(args *CityShowArgs, reply *CityShowReply) error {
	v, err := s.w.CityShow(args.UserId, args.CharacterId, args.CityId)
	reply.View = v
	return err
}

func (s *RegionService) CityStudy(args *CityStudyArgs, reply *CityStudyReply) error {
	id, err := s.w.CityStudy(args.UserId, args.CharacterId, args.CityId, args.KnowledgeId)
	reply.Id = id
	return err
}

func (s *RegionService) CityBuild(args *CityBuildArgs, reply *CityBuildReply) error {
	id, err := s.w.CityBuild(args.UserId, args.CharacterId, args.CityId, args.BuildingId)
	reply.Id = id
	return err
}

func (s *RegionService) CityTrain(args *CityTrainArgs, reply *CityTrainReply) error {
	id, err := s.w.CityTrain(args.UserId, args.CharacterId, args.CityId, args.UnitId)
	reply.Id = id
	return err
}

func (s *RegionService) MapDot(args *MapDotArgs, reply *MapDotReply) error {
	reply.Dot = s.w.Places.Dot()
	return nil
}

func (s *RegionService) MapCheck(args *MapCheckArgs, reply *MapCheckReply) error {
	s.w.Places.Rehash()
	return nil
}

func (s *RegionService) MapRehash(args *MapRehashArgs, reply *MapRehashReply) error {
	return s.w.Places.Check(s.w)
}

func (s *RegionService) MapPlaces(args *MapPlacesArgs, reply *MapPlacesReply) error {
	reply.Items = s.w.Places
	return nil
}

func (s *RegionService) MapCities(args *MapCitiesArgs, reply *MapCitiesReply) error {
	reply.Items = s.w.Live.Cities
	return nil
}

func (s *RegionService) MapArmies(args *MapArmiesArgs, reply *MapArmiesReply) error {
	reply.Items = s.w.Live.Armies
	return nil
}

func (s *RegionService) AdminSave(args *AdminSaveArgs, reply *AdminSaveReply) error {
	return save(s.w)
}

func (s *RegionService) AdminCheck(args *AdminCheckArgs, reply *AdminCheckReply) error {
	return s.w.Check()
}

func (s *RegionService) RoundProduce(args *RoundProduceArgs, reply *RoundProduceReply) error {
	s.w.Produce()
	return nil
}

func (s *RegionService) RoundMove(args *RoundMoveArgs, reply *RoundMoveReply) error {
	s.w.Move()
	return nil
}

func main() {
	var err error
	var w World

	w.Init()

	var north string
	var pathLoad string
	flag.StringVar(&north, "north", "127.0.0.1:8081", "File to be loaded")
	flag.StringVar(&pathLoad, "load", "", "File to be loaded")
	flag.StringVar(&pathSave, "save", "/tmp/hegemonie/data", "Directory for persistent")
	flag.Parse()

	if pathSave != "" {
		err = os.MkdirAll(pathSave, 0755)
		if err != nil {
			log.Fatalf("Failed to create [%s]: %s", pathSave, err.Error())
		}
	}

	if pathLoad != "" {
		type cfgSection struct {
			suffix string
			obj    interface{}
		}
		cfgSections := []cfgSection{
			{"defs.json", &w.Definitions},
			{"map.json", &w.Places},
			{"auth.json", &w.Auth},
			{"live.json", &w.Live},
		}
		for _, section := range cfgSections {
			var in io.ReadCloser
			p := pathLoad + "/" + section.suffix
			in, err = os.Open(p)
			if err != nil {
				log.Fatalf("Failed to load the World from [%s]: %s", p, err.Error())
			}
			err = json.NewDecoder(in).Decode(section.obj)
			in.Close()
			if err != nil {
				log.Fatalf("Failed to load the World from [%s]: %s", p, err.Error())
			}
		}
		err = w.PostLoad()
		if err != nil {
			log.Fatalf("Inconsistent World from [%s]: %s", pathLoad, err.Error())
		}
	}

	err = w.Check()
	if err != nil {
		log.Fatalf("Inconsistent World: %s", err.Error())
	}

	var srv Region = &RegionService{w: &w}

	err = rpc.RegisterName("Region", srv)
	if err != nil {
		log.Fatalf("RPC error: %s", err.Error())
	}

	rpc.HandleHTTP()
	err = http.ListenAndServe(north, nil)

	if err != nil {
		log.Printf("Server error: %s", err.Error())
	}

	if pathSave != "" {
		err = save(&w)
		if err != nil {
			log.Fatalf("Failed to save the World at exit: %s", err.Error())
		}
	}
}
