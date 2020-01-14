// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/google/subcommands"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	. "github.com/jfsmig/hegemonie/common/client"
	"github.com/jfsmig/hegemonie/common/world"
)

type RegionCommand struct {
	north    string
	pathLoad string
	srv      RegionService
}

type RegionService struct {
	w        *world.World
	pathSave string
}

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}

func (self *RegionService) save() error {
	if self.pathSave == "" {
		return errors.New("No save path configured")
	}
	p := self.pathSave + "/" + makeSaveFilename()
	p = filepath.Clean(p)
	out, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = self.w.DumpJSON(out)
	out.Close()
	if err != nil {
		_ = os.Remove(p)
		return err
	}

	latest := self.pathSave + "/latest"
	_ = os.Remove(latest)
	_ = os.Symlink(p, latest)
	return nil
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

func (s *RegionService) CityCreateArmy(args *CityCreateArmyArgs, reply *CityCreateArmyReply) error {
	id, err := s.w.CityCreateArmy(args.UserId, args.CharacterId, args.CityId, args.Name)
	reply.Id = id
	return err
}

func (s *RegionService) CityTransferUnit(args *CityTransferUnitArgs, reply *CityTransferUnitReply) error {
	return s.w.CityTransferUnit(args.UserId, args.CharacterId, args.CityId, args.UnitId, args.ArmyId)
}

func (s *RegionService) CityCommandArmy(args *CityCommandArmyArgs, reply *CityCommandArmyReply) error {
	//return s.w.CityCommandArmy(args.UserId, args.CharacterId, args.CityId, args.UnitId, args.ArmyId)
	return errors.New("NYI")
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
	return s.save()
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

func (s *RegionService) GetScoreBoard(args *GetScoreBoardArgs, reply *GetScoreBoardReply) error {
	reply.Board = s.w.ScoreBoardCompute()
	return nil
}

func (self *RegionCommand) Name() string     { return "region" }
func (self *RegionCommand) Synopsis() string { return "Start a region service." }
func (self *RegionCommand) Usage() string    { return "region ENDPOINT\n" }

func (self *RegionCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&self.north, "north", ":8080", "File to be loaded")
	f.StringVar(&self.pathLoad, "load", "/data/defs", "File to be loaded")
	f.StringVar(&self.srv.pathSave, "save", "/data/dump", "Directory for persistent")
}

func (self *RegionCommand) Execute(_ context.Context, f0 *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error

	self.srv.w = new(world.World)
	self.srv.w.Init()

	if self.srv.pathSave != "" {
		err = os.MkdirAll(self.srv.pathSave, 0755)
		if err != nil {
			log.Printf("Failed to create [%s]: %s", self.srv.pathSave, err.Error())
			return subcommands.ExitFailure
		}
	}

	if self.pathLoad != "" {
		type cfgSection struct {
			suffix string
			obj    interface{}
		}
		cfgSections := []cfgSection{
			{"defs.json", &self.srv.w.Definitions},
			{"map.json", &self.srv.w.Places},
			{"auth.json", &self.srv.w.Auth},
			{"live.json", &self.srv.w.Live},
		}
		for _, section := range cfgSections {
			var in io.ReadCloser
			p := self.pathLoad + "/" + section.suffix
			in, err = os.Open(p)
			if err != nil {
				log.Printf("Failed to load the World from [%s]: %s", p, err.Error())
				return subcommands.ExitFailure
			}
			err = json.NewDecoder(in).Decode(section.obj)
			in.Close()
			if err != nil {
				log.Printf("Failed to load the World from [%s]: %s", p, err.Error())
				return subcommands.ExitFailure
			}
		}
		err = self.srv.w.PostLoad()
		if err != nil {
			log.Printf("Inconsistent World from [%s]: %s", self.pathLoad, err.Error())
			return subcommands.ExitFailure
		}
	}

	err = self.srv.w.Check()
	if err != nil {
		log.Printf("Inconsistent World: %s", err.Error())
		return subcommands.ExitFailure
	}

	var srv Region = &RegionService{w: self.srv.w}

	err = rpc.RegisterName("Region", srv)
	if err != nil {
		log.Printf("RPC error: %s", err.Error())
		return subcommands.ExitFailure
	}

	rpc.HandleHTTP()
	err = http.ListenAndServe(self.north, nil)

	if err != nil {
		log.Printf("Server error: %s", err.Error())
		return subcommands.ExitFailure
	}

	if self.srv.pathSave != "" {
		err = self.srv.save()
		if err != nil {
			log.Printf("Failed to save the World at exit: %s", err.Error())
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}
