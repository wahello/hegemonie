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

func (s *srvCity) List(req *proto.CitiesByCharReq, stream proto.City_ListServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	r := s.w.Regions.Get(req.Region)
	if r == nil {
		return status.Error(codes.NotFound, "No such region")
	}

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
}

func (s *srvCity) AllCities(req *proto.PaginatedQuery, stream proto.City_AllCitiesServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	r := s.w.Regions.Get(req.Region)
	if r == nil {
		return status.Error(codes.NotFound, "No such region")
	}

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
}

func (s *srvCity) Show(ctx context.Context, req *proto.CityId) (*proto.CityView, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	r := s.w.Regions.Get(req.Region)
	if r == nil {
		return nil, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such city")
	}

	view := ShowCity(s.w, city)
	utils.Logger.Debug().
		Int("#a", len(view.Assets.Armies)).
		Int("#k", len(view.Assets.Knowledges)).
		Int("#b", len(view.Assets.Buildings)).
		Int("#u", len(view.Assets.Units)).
		Interface("prod", view.Production).
		Interface("stock", view.Stock).
		Msg("city")
	return view, nil
}

func (s *srvCity) Study(ctx context.Context, req *proto.StudyReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return nil, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such city")
	}

	_, err = city.Study(r, req.KnowledgeType)
	return none, nil
}

func (s *srvCity) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return nil, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such city")
	}

	_, err = city.Build(r, req.BuildingType)
	return none, nil
}

func (s *srvCity) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return nil, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such city")
	}

	_, err = city.Train(r, req.UnitType)
	return none, nil
}

func (s *srvCity) ListArmies(req *proto.CityId, stream proto.City_ListArmiesServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	r := s.w.Regions.Get(req.GetRegion())
	if r == nil {
		return status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCharacter(), req.GetCity())
	if err != nil {
		return status.Error(codes.NotFound, "No such city")
	}

	var last string
	for {
		tab := city.Armies.Slice(last, 100)
		if len(tab) <= 0 {
			return nil
		}
		for _, a := range city.Armies {
			last = a.ID
			err = stream.Send(&proto.ArmyName{Id: a.ID, Name: a.Name})
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

// Create an army made of only Units (no Resources carried)
func (s *srvCity) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return none, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return none, status.Error(codes.NotFound, "No such city")
	}

	_, err = city.CreateArmyFromIds(r, req.Unit...)
	return none, err
}

// Create an army made of only Resources (no Units)
func (s *srvCity) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return none, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return none, status.Error(codes.NotFound, "No such city")
	}

	resources := resAbsP2M(req.Stock)
	_, err = city.CreateTransport(r, resources)
	return none, err
}

func (s *srvCity) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return none, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return none, status.Error(codes.NotFound, "No such city")
	}

	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, status.Error(codes.NotFound, "No such army")
	}

	err = city.TransferOwnUnit(army, req.Unit...)
	return none, err
}

func (s *srvCity) TransferResources(ctx context.Context, req *proto.TransferResourcesReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	r := s.w.Regions.Get(req.GetCity().GetRegion())
	if r == nil {
		return none, status.Error(codes.NotFound, "No such region")
	}

	city, err := r.CityGetAndCheck(req.GetCity().GetCharacter(), req.GetCity().GetCity())
	if err != nil {
		return none, status.Error(codes.NotFound, "No such city")
	}

	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, status.Error(codes.NotFound, "No such army")
	}

	err = city.TransferOwnResources(army, resAbsP2M(req.Stock))
	return none, err
}
