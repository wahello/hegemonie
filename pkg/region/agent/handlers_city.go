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
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	armies, err := s.w.CityArmies(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	rep := &proto.ListOfNamedItems{}
	for _, a := range armies {
		rep.Items = append(rep.Items, &proto.NamedItem{Id: a.Id, Name: a.Name})
	}
	return rep, nil
}

func (s *srvCity) List(ctx context.Context, req *proto.ListReq) (*proto.ListOfNamedItems, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	rep := &proto.ListOfNamedItems{}
	cities := s.w.Cities(req.Character)

	for _, c := range cities {
		rep.Items = append(rep.Items, &proto.NamedItem{Id: c.Id, Name: c.Name})
	}
	return rep, nil
}

func (s *srvCity) Show(ctx context.Context, req *proto.CityId) (*proto.CityView, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	view := ShowCity(s.w, city)
	return view, nil
}

func (s *srvCity) Study(ctx context.Context, req *proto.StudyReq) (*proto.None, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Study(s.w, req.KnowledgeType)
	return &proto.None{}, err
}

func (s *srvCity) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Build(s.w, req.BuildingType)
	return &proto.None{}, err
}

func (s *srvCity) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}

	_, err = city.Train(s.w, req.UnitType)
	return &proto.None{}, err
}

func (s *srvCity) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvCity) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvCity) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}
