// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type srvAdmin struct {
	cfg *Config
	w   *region.World
}

var none = &proto.None{}

func (sa *srvAdmin) rlockDo(action func() error) error {
	sa.w.RLock()
	defer sa.w.RUnlock()
	return action()
}

func (sa *srvAdmin) wlockDo(action func() error) error {
	sa.w.WLock()
	defer sa.w.WUnlock()
	return action()
}

func (sa *srvAdmin) Produce(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, sa.rlockDo(func() error {
		r := sa.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		r.Produce()
		return nil
	})
}

func (sa *srvAdmin) Move(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, sa.rlockDo(func() error {
		r := sa.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		r.Move()
		return nil
	})
}

func (sa *srvAdmin) CreateRegion(ctx context.Context, req *proto.RegionCreateReq) (*proto.None, error) {
	return none, sa.wlockDo(func() error {
		_, err := sa.w.CreateRegion(req.Name, req.MapName)
		return err
	})
}

func (sa *srvAdmin) GetScores(req *proto.RegionId, stream proto.Admin_GetScoresServer) error {
	return sa.rlockDo(func() error {
		r := sa.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		for _, c := range r.Cities {
			err := stream.Send(ShowCityPublic(sa.w, c, true))
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}
