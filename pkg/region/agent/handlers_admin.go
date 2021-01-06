// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/discovery"
	mproto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"google.golang.org/grpc"
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
	//  first, load the cities from the maps repository
	endpoint, err := discovery.DefaultDiscovery.Map()
	if err != nil {
		return none, err
	}
	cnx, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return none, err
	}
	defer cnx.Close()

	client := mproto.NewMapClient(cnx)

	marker := uint64(0)
	rep, err := client.Cities(ctx, &mproto.ListCitiesReq{
		MapName: req.MapName,
		Marker:  marker,
	})
	if err != nil {
		return none, err
	}

	out := make([]region.NamedCity, 0)
	for {
		x, err := rep.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return none, err
		}
		marker = x.GetId()
		out = append(out, region.NamedCity{Name: x.GetName(), ID: x.GetId()})
	}

	return none, sa.wlockDo(func() error {
		_, err := sa.w.CreateRegion(req.Name, req.MapName, out)
		return err
	})
}

func (sa *srvAdmin) ListRegions(req *proto.RegionListReq, stream proto.Admin_ListRegionsServer) error {
	marker := req.NameMarker
	return sa.rlockDo(func() error {
		for {
			tab := sa.w.Regions.Slice(marker, 100)
			if len(tab) <= 0 {
				break
			}
			for _, x := range tab {
				marker = x.Name
				summary := &proto.RegionSummary{
					Name:        x.Name,
					MapName:     x.MapName,
					CountCities: uint32(len(x.Cities)),
					CountFights: uint32(len(x.Fights)),
				}
				err := stream.Send(summary)
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
			}
		}
		return nil
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
