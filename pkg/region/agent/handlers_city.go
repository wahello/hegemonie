// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvCity struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvCity) List(ctx context.Context, req *proto.ListReq) (*proto.ListOfCities, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	rep := &proto.ListOfCities{}
	cities := s.w.Cities(req.Character)

	for _, c := range cities {
		rep.Items = append(rep.Items, ShowCityPublic(s.w, c, false))
	}
	return rep, nil
}

func (s *srvCity) Show(ctx context.Context, req *proto.CityId) (*proto.CityView, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
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

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Study(s.w, req.KnowledgeType)
	return &proto.None{}, err
}

func (s *srvCity) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Build(s.w, req.BuildingType)
	return &proto.None{}, err
}

func (s *srvCity) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Train(s.w, req.UnitType)
	return &proto.None{}, err
}

func (s *srvCity) ListArmies(ctx context.Context, req *proto.CityId) (*proto.ListOfNamedItems, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	rep := &proto.ListOfNamedItems{}
	for _, a := range city.Armies {
		rep.Items = append(rep.Items, &proto.NamedItem{Id: a.ID, Name: a.Name})
	}
	return rep, nil
}

// Create an army made of only Units (no Resources carried)
func (s *srvCity) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.CreateArmyFromIds(s.w, req.Unit...)
	return &proto.None{}, err
}

// Create an army made of only Resources (no Units)
func (s *srvCity) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	r := resAbsP2M(req.Stock)
	_, err = city.CreateTransport(s.w, r)
	return &proto.None{}, err
}

func (s *srvCity) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "Army not found (id %v)", req.Army)
	}

	err = city.TransferOwnUnit(army, req.Unit...)
	return &proto.None{}, err
}

func (s *srvCity) TransferResources(ctx context.Context, req *proto.TransferResourcesReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "City Not found")
	}

	err = city.TransferOwnResources(army, resAbsP2M(req.Stock))
	return &proto.None{}, err
}
