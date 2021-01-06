// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"encoding/json"
	"github.com/jfsmig/hegemonie/pkg/discovery"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"google.golang.org/grpc"
	"io"
	"os"
)

type ClientCLI struct{}

func (cli *ClientCLI) DoCreateRegion(ctx context.Context, args []string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewAdminClient(cnx)
		_, err := client.CreateRegion(ctx, &proto.RegionCreateReq{MapName: args[1], Name: args[0]})
		if err != nil {
			return err
		}

		type status struct {
			Msg  string
			Code int
			ID   string
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(status{Msg: "Create", Code: 200, ID: args[0]})
		return nil

	})
}

func (cli *ClientCLI) DoListRegions(ctx context.Context, args []string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		marker := ""
		if len(args) > 0 {
			marker = args[0]
		}
		client := proto.NewAdminClient(cnx)
		rep, err := client.ListRegions(ctx, &proto.RegionListReq{NameMarker: marker})
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "")
		for {
			x, err := rep.Recv()
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
			err = encoder.Encode(x)
			if err != nil {
				panic(err)
			}
		}
		return nil

	})
}

func (cli *ClientCLI) DoRegionMovement(ctx context.Context, args []string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "")
		client := proto.NewAdminClient(cnx)
		for _, reg := range args {
			rep, err := client.Move(ctx, &proto.RegionId{Region: reg})
			if err != nil {
				return err
			}
			err = encoder.Encode(rep)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
}

func (cli *ClientCLI) DoRegionProduction(ctx context.Context, args []string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "")
		client := proto.NewAdminClient(cnx)
		for _, reg := range args {
			rep, err := client.Produce(ctx, &proto.RegionId{Region: reg})
			if err != nil {
				return err
			}
			err = encoder.Encode(rep)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
}

type actionFunc func(ctx context.Context, cli *grpc.ClientConn) error

func (cli *ClientCLI) connect(ctx context.Context, action actionFunc) error {
	endpoint, err := discovery.DefaultDiscovery.Region()
	if err != nil {
		return err
	}
	cnx, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	return action(ctx, cnx)
}
