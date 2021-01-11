// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"

	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvCity struct {
	cfg *Config
	w   *region.World
}

func (s *srvCity) _regLock(mode rune, regID string, action func(*region.Region) error) error {
	switch mode {
	case 'r':
		s.w.RLock()
		defer s.w.RUnlock()
	case 'w':
		s.w.WLock()
		defer s.w.WUnlock()
	default:
		panic("wtf!?")
	}
	r := s.w.Regions.Get(regID)
	if r == nil {
		return status.Error(codes.NotFound, "no such region")
	}
	return action(r)
}

func (s *srvCity) _cityLock(mode rune, regID, charID string, cityID uint64, action func(*region.Region, *region.City) error) error {
	return s._regLock('r', regID, func(r *region.Region) error {
		switch mode {
		case 'r':
			s.w.RLock()
			defer s.w.RUnlock()
		case 'w':
			s.w.WLock()
			defer s.w.WUnlock()
		default:
			panic("wtf!?")
		}
		city, e := r.CityGetAndCheck(charID, cityID)
		if e != nil {
			return status.Error(codes.NotFound, "no such city")
		}
		return action(r, city)
	})
}

func (s *srvCity) cityLock(mode rune, req *proto.CityId, action func(*region.Region, *region.City) error) error {
	return s._cityLock(mode, req.Region, req.Character, req.City, action)
}

func (s *srvCity) List(req *proto.CitiesByCharReq, stream proto.City_ListServer) error {
	return s._regLock('r', req.Region, func(r *region.Region) error {
		last := req.Marker
		for {
			tab := r.Cities.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, c := range tab {
				last = c.ID
				if c.Owner != req.Character && c.Deputy != req.Character {
					continue
				}
				err := stream.Send(ShowCityPublic(s.w, c, false))
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

func (s *srvCity) AllCities(req *proto.PaginatedQuery, stream proto.City_AllCitiesServer) error {
	return s._regLock('r', req.Region, func(r *region.Region) error {
		last := req.Marker
		for {
			tab := r.Cities.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, c := range tab {
				last = c.ID
				err := stream.Send(ShowCityPublic(s.w, c, false))
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

func (s *srvCity) Show(ctx context.Context, req *proto.CityId) (reply *proto.CityView, err error) {
	err = s.cityLock('r', req, func(r *region.Region, c *region.City) error {
		view := ShowCity(s.w, c)
		utils.Logger.Debug().
			Int("#a", len(view.Assets.Armies)).
			Int("#k", len(view.Assets.Knowledges)).
			Int("#b", len(view.Assets.Buildings)).
			Int("#u", len(view.Assets.Units)).
			Interface("prod", view.Production).
			Interface("stock", view.Stock).
			Msg("city")
		reply = view
		return nil
	})
	return reply, err
}

func (s *srvCity) Study(ctx context.Context, req *proto.StudyReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Study(r, req.KnowledgeType)
		return e
	})
}

func (s *srvCity) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Build(r, req.BuildingType)
		return e
	})
}

func (s *srvCity) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Train(r, req.UnitType)
		return e
	})
}

func (s *srvCity) ListArmies(req *proto.CityId, stream proto.City_ListArmiesServer) error {
	return s.cityLock('r', req, func(r *region.Region, c *region.City) error {
		var last string
		for {
			tab := c.Armies.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, a := range c.Armies {
				last = a.ID
				err := stream.Send(&proto.ArmyName{Id: a.ID, Name: a.Name})
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

// Create an army made of only Units (no Resources carried)
func (s *srvCity) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.CreateArmyFromIds(r, req.Unit...)
		return e
	})
}

// Create an army made of only Resources (no Units)
func (s *srvCity) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.CreateTransport(r, resAbsP2M(req.Stock))
		return e
	})
}

func (s *srvCity) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		army := c.Armies.Get(req.Army)
		if army == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return c.TransferOwnUnit(army, req.Unit...)
	})
}

func (s *srvCity) TransferResources(ctx context.Context, req *proto.TransferResourcesReq) (*proto.None, error) {
	return none, s.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		army := c.Armies.Get(req.Army)
		if army == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return c.TransferOwnResources(army, resAbsP2M(req.Stock))
	})
}
