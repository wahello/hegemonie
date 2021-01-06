// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type srvArmy struct {
	cfg *Config
	w   *region.World
}

func (s *srvArmy) getAndCheckArmy(req *proto.ArmyId) (*region.Region, *region.City, *region.Army, error) {
	r := s.w.Regions.Get(req.Region)
	if r == nil {
		return nil, nil, nil, status.Error(codes.NotFound, "no such region")
	}
	city, err := r.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, nil, nil, status.Error(codes.NotFound, "no such city")
	}
	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, nil, nil, status.Error(codes.NotFound, "no such army")
	}
	return r, city, army, err
}

func (s *srvArmy) wlock(action func() error) error {
	s.w.WLock()
	defer s.w.WUnlock()
	return action()
}

func (s *srvArmy) rlock(action func() error) error {
	s.w.RLock()
	defer s.w.RUnlock()
	return action()
}

type actionFunc func(*region.Region, *region.City, *region.Army) error

func (s *srvArmy) getAndDo(id *proto.ArmyId, action actionFunc) error {
	r, city, army, err := s.getAndCheckArmy(id)
	if err == nil {
		err = action(r, city, army)
	}
	return err
}

func (s *srvArmy) wlockDo(id *proto.ArmyId, action actionFunc) error {
	return s.wlock(func() error { return s.getAndDo(id, action) })
}

func (s *srvArmy) rlockDo(id *proto.ArmyId, action actionFunc) error {
	return s.rlock(func() error { return s.getAndDo(id, action) })
}

func (s *srvArmy) Show(ctx context.Context, req *proto.ArmyId) (*proto.ArmyView, error) {
	var rc *proto.ArmyView
	err := s.rlockDo(req, func(_ *region.Region, _ *region.City, army *region.Army) error {
		rc = ShowArmy(s.w, army)
		return nil
	})
	return rc, err
}

func (s *srvArmy) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.Flea(r)
	})
}

func (s *srvArmy) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.Flip(r)
	})
}

func (s *srvArmy) Move(ctx context.Context, req *proto.ArmyMoveReq) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.DeferMove(r, req.Target, region.ActionArgMove{})
		})
}

func (s *srvArmy) Attack(ctx context.Context, req *proto.ArmyAssaultReq) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.DeferAttack(r, req.Target, region.ActionArgAssault{})
		})
}

func (s *srvArmy) Wait(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.DeferWait(r, req.Target)
		})
}

func (s *srvArmy) Defend(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.DeferDefend(r, req.Target)
		})
}

func (s *srvArmy) Disband(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.DeferDisband(r, req.Target)
		})
}

func (s *srvArmy) Cancel(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req,
		func(r *region.Region, _ *region.City, army *region.Army) error {
			return army.Cancel(r)
		})
}
