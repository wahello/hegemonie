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

type srvArmy struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvArmy) getAndCheckArmy(req *proto.ArmyId) (*region.City, *region.Army, error) {
	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, nil, err
	}
	if army := s.w.ArmyGet(req.Army); army == nil {
		return nil, nil, status.Errorf(codes.NotFound, "Army Not found")
	} else if army.City != city.Id {
		return nil, nil, status.Errorf(codes.PermissionDenied, "Army not controlled")
	} else {
		return city, army, err
	}
}

func (s *srvArmy) Show(ctx context.Context, req *proto.ArmyId) (*proto.ArmyView, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	if _, army, err := s.getAndCheckArmy(req); err != nil {
		return nil, err
	} else {
		return ShowArmy(s.w, army), nil
	}
}

func (s *srvArmy) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	if _, army, err := s.getAndCheckArmy(req); err != nil {
		return nil, err
	} else {
		if err = army.Flea(s.w); err != nil {
			return nil, err
		} else {
			return &proto.None{}, nil
		}
	}
}

func (s *srvArmy) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	if _, army, err := s.getAndCheckArmy(req); err != nil {
		return nil, err
	} else {
		if err = army.Flip(s.w); err != nil {
			return nil, err
		} else {
			return &proto.None{}, nil
		}
	}
}

func (s *srvArmy) Command(ctx context.Context, req *proto.ArmyCommandReq) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	_, army, err := s.getAndCheckArmy(req.Id)
	if err != nil {
		return nil, err
	}
	target := s.w.CityGet(req.Command.Target)
	if target == nil {
		return nil, status.Errorf(codes.NotFound, "Target Not found")
	}

	switch req.Command.Action {
	case region.CmdCityAttack:
		err = army.DeferAttack(s.w, target)
	case region.CmdCityDefend:
		err = army.DeferDefend(s.w, target)
	case region.CmdCityOverlord:
		err = army.DeferConquer(s.w, target)
	case region.CmdCityLiberate:
		err = army.DeferLiberate(s.w, target)
	case region.CmdCityBreak:
		err = army.DeferBreak(s.w, target)
	case region.CmdCityMassacre:
		err = army.DeferMassacre(s.w, target)
	case region.CmdCityDeposit:
		err = army.DeferDeposit(s.w, target)
	case region.CmdCityDisband:
		err = army.DeferDisband(s.w, target)
	default:
		return nil, status.Errorf(codes.NotFound, "Invalid action")
	}

	if err != nil {
		return nil, err
	} else {
		return &proto.None{}, nil
	}
}
