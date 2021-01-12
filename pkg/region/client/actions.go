// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

type ClientCLI struct{}

// DoCreateRegion triggers the synchronous creation of a region with the given name, modeled on the named map.
func (cli *ClientCLI) DoCreateRegion(ctx context.Context, regID, mapID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).CreateRegion(ctx, &proto.RegionCreateReq{MapName: mapID, Name: regID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, regID, "created")
	})
}

// DoListRegions dumps to os.Stdout a JSON stream of the known regions, sorted by name
func (cli *ClientCLI) DoListRegions(ctx context.Context, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewAdminClient(cnx).ListRegions(ctx, &proto.RegionListReq{NameMarker: marker})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// DoRegionMovement triggers one round of armies movement on all the cities of the named region
func (cli *ClientCLI) DoRegionMovement(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Move(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Moved")
	})
}

// DoRegionProduction triggers the production of resources on all the cities of the named region
func (cli *ClientCLI) DoRegionProduction(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Produce(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Produced")
	})
}

func (cli *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Region()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}
