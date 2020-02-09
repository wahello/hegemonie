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
	"log"

	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvArmy struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvArmy) Show(ctx context.Context, req *proto.ArmyId) (*proto.ArmyView, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	army := s.w.ArmyGet(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "Army Not found")
	}
	log.Println(army)
	log.Println(city)
	if army.City != city.Id {
		return nil, status.Errorf(codes.PermissionDenied, "Army not controlled")
	}

	return ShowArmy(s.w, army), nil
}

func (s *srvArmy) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	army := s.w.ArmyGet(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "Army Not found")
	}
	if army.City != city.Id {
		return nil, status.Errorf(codes.PermissionDenied, "Army not controlled")
	}

	if err = army.Flea(s.w); err != nil {
		return nil, err
	} else {
		return &proto.None{}, nil
	}
}

func (s *srvArmy) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	army := s.w.ArmyGet(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "Army Not found")
	}
	if army.City != city.Id {
		return nil, status.Errorf(codes.PermissionDenied, "Army not controlled")
	}

	if err = army.Flip(s.w); err != nil {
		return nil, err
	} else {
		return &proto.None{}, nil
	}
}

func (s *srvArmy) Command(ctx context.Context, req *proto.ArmyCommandReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, err
	}
	army := s.w.ArmyGet(req.Army)
	if army == nil {
		return nil, status.Errorf(codes.NotFound, "Army Not found")
	}
	if army.City != city.Id {
		return nil, status.Errorf(codes.PermissionDenied, "Army not controlled")
	}

	// FIXME(jfs): NYI
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}
