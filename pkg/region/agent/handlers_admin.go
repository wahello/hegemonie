// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	mproto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type adminApp struct {
	regionApp
}

func (app *regionApp) Produce(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, app._regLock('w', req.Region, func(r *region.Region) error {
		r.Produce(ctx)
		return nil
	})
}

func (app *regionApp) Move(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, app._regLock('w', req.Region, func(r *region.Region) error {
		r.Move(ctx)
		return nil
	})
}

func (app *regionApp) CreateRegion(ctx context.Context, req *proto.RegionCreateReq) (*proto.None, error) {
	//  first, load the cities from the maps repository
	endpoint, err := utils.DefaultDiscovery.Map()
	if err != nil {
		return none, status.Errorf(codes.Internal, "configuration error: %v", err)
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

	return none, app._worldLock('w', func() error {
		_, err := app.w.CreateRegion(req.Name, req.MapName, out)
		return err
	})
}

func (app *regionApp) ListRegions(req *proto.RegionListReq, stream proto.Admin_ListRegionsServer) error {
	return app._worldLock('r', func() error {
		for marker := req.NameMarker; ; {
			tab := app.w.Regions.Slice(marker, 100)
			if len(tab) <= 0 {
				return nil
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
						return nil
					}
					return err
				}
			}
		}
	})
}

func (app *regionApp) GetScores(req *proto.RegionId, stream proto.Admin_GetScoresServer) error {
	return app._regLock('r', req.Region, func(r *region.Region) error {
		for _, c := range r.Cities {
			err := stream.Send(ShowCityPublic(app.w, c, true))
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
