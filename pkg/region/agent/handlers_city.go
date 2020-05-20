// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvCity struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvCity) ListArmies(ctx context.Context, req *proto.CityId) (*proto.ListOfNamedItems, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	rep := &proto.ListOfNamedItems{}
	for _, a := range city.Armies() {
		rep.Items = append(rep.Items, &proto.NamedItem{Id: a.Id, Name: a.Name})
	}
	return rep, nil
}

func (s *srvCity) List(ctx context.Context, req *proto.ListReq) (*proto.ListOfCities, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	rep := &proto.ListOfCities{}
	cities := s.w.Cities(req.Character)

	for _, c := range cities {
		rep.Items = append(rep.Items, &proto.PublicCity{
			Id: c.Id, Name: c.Name, Cell: c.Cell,
			Politics:  c.PoliticalGroup,
			Chaos:     c.Chaotic,
			Alignment: c.Alignment,
			Ethny:     c.EthnicGroup,
			Cult:      c.Cult,
		})
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
	return view, nil
}

func (s *srvCity) Study(ctx context.Context, req *proto.StudyReq) (*proto.None, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Study(s.w, req.KnowledgeType)
	return &proto.None{}, err
}

func (s *srvCity) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Build(s.w, req.BuildingType)
	return &proto.None{}, err
}

func (s *srvCity) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Train(s.w, req.UnitType)
	return &proto.None{}, err
}

func (s *srvCity) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	for _, uid := range req.Unit {
		if city.Unit(uid) == nil {
			return nil, status.Errorf(codes.NotFound, "Troop not found (id %v)", uid)
		}
	}

	army, err := s.w.ArmyCreate(city, req.Name)
	if err != nil {
		return nil, err
	}

	for _, uid := range req.Unit {
		city.TransferOwnUnit(army, uid)
	}
	return &proto.None{}, nil
}

func (s *srvCity) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	army := s.w.ArmyGet(req.Army)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Army not found (id %v)", req.Army)
	}

	err = city.TransferOwnUnit(army, req.Unit...)
	if err != nil {
		return nil, err
	}

	return &proto.None{}, nil
}

func (s *srvCity) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	r := resAbsP2M(req.Stock)
	if !city.Stock.GreaterOrEqualTo(r) {
		return nil, status.Errorf(codes.FailedPrecondition, "Insufficient resources")
	}
	army, err := s.w.ArmyCreate(city, req.Name)
	city.Stock.Remove(r)
	army.Stock.Add(r)
	return &proto.None{}, nil
}

func (s *srvCity) TransferResources(ctx context.Context, req *proto.TransferResourcesReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	army := s.w.ArmyGet(req.Army)
	if army != nil {
		return nil, status.Errorf(codes.NotFound, "City Not found")
	}

	err = city.TransferOwnResources(army, resAbsP2M(req.Stock))
	if err != nil {
		return nil, err
	}
	return &proto.None{}, nil
}
