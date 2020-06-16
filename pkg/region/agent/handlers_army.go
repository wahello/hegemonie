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
	if army := city.Armies.Get(req.Army); army == nil {
		return nil, nil, status.Errorf(codes.NotFound, "Army Not found")
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
	target := s.w.Places.CellGet(req.Command.Target)
	if target == nil {
		err = status.Errorf(codes.NotFound, "Target Not found")
	} else {
		switch req.Command.Action {
		case proto.ArmyCommandType_Move:
			err = army.DeferMove(s.w, target)
		case proto.ArmyCommandType_Attack:
			err = army.DeferAttack(s.w, target)
		case proto.ArmyCommandType_Defend:
			err = army.DeferDefend(s.w, target)
		case proto.ArmyCommandType_Wait:
			err = army.DeferWait(s.w, target)
		case proto.ArmyCommandType_Overlord:
			err = army.DeferOverlord(s.w, target)
		case proto.ArmyCommandType_Break:
			err = army.DeferBreak(s.w, target)
		case proto.ArmyCommandType_Massacre:
			err = army.DeferMassacre(s.w, target)
		case proto.ArmyCommandType_Deposit:
			err = army.DeferDeposit(s.w, target)
		case proto.ArmyCommandType_Disband:
			err = army.DeferDisband(s.w, target)
		default:
			err = status.Errorf(codes.InvalidArgument, "Invalid action")
		}
	}

	if err != nil {
		return nil, err
	} else {
		return &proto.None{}, nil
	}
}
