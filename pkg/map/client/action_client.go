// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jfsmig/hegemonie/pkg/discovery"
	"github.com/jfsmig/hegemonie/pkg/map/proto"
	"google.golang.org/grpc"
	"io"
	"os"
	"strconv"
)

// ClientCLI gathers the command line interface actions related to the Map API service.
// It wraps the main gRPC client calls and dumps all the output as JSON on os.Stdout.
type ClientCLI struct{}

// GetCities produces to os.Stdout a JSON array of cities, where each City is an <id,name> tuple.
func (c *ClientCLI) GetCities(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewMapClient(cnx)

		rep, err := client.Cities(ctx, &proto.ListCitiesReq{
			MapName: args.MapName,
			Marker:  args.Src,
		})
		if err != nil {
			return err
		}

		out := make([]uint64, 0)
		for {
			x, err := rep.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			out = append(out, x.GetId())
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)
		return nil
	})
}

// GetRoads produces to os.Stdout a JSON array of <int,int> pairs, with one pair
// for a road in place on the given map.
func (c *ClientCLI) GetRoads(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewMapClient(cnx)

		rep, err := client.Edges(ctx, &proto.ListEdgesReq{
			MapName:   args.MapName,
			MarkerSrc: args.Src,
			MarkerDst: args.Dst,
		})
		if err != nil {
			return err
		}

		type Pair struct{ Src, Dst uint64 }
		out := make([]Pair, 0)
		for {
			x, err := rep.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			out = append(out, Pair{x.GetSrc(), x.GetDst()})
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)
		return nil
	})
}

// GetPositions produces to os.Stdout a JSON array of <id,x,y> tuples, all fields being integers,
// with one tuple for each position on the map, i.e. one tuple per vertex in the graph.
func (c *ClientCLI) GetPositions(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewMapClient(cnx)

		rep, err := client.Vertices(ctx, &proto.ListVerticesReq{
			MapName: args.MapName,
			Marker:  args.Src,
		})
		if err != nil {
			return err
		}

		type V struct{ ID, X, Y uint64 }
		out := make([]V, 0)
		for {
			x, err := rep.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			out = append(out, V{ID: x.GetId(), X: x.GetX(), Y: x.GetY()})
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)
		return nil
	})
}

// GetPath produces to os.Stdout a JSON array of integers, all fields being the unique ID
// of a position involved in the path from the given source to the given destination.
func (c *ClientCLI) GetPath(ctx context.Context, args PathArgs) error {
	return c.getPath(ctx, args)
}

// GetStep produces to os.Stdout a JSON singleton array of integers, with the unique value
// equals to the ID of the position that is the next step in the path from the given source
// to the given destination.
func (c *ClientCLI) GetStep(ctx context.Context, args PathArgs) error {
	args.Max = 1
	return c.getPath(ctx, args)
}

func (c *ClientCLI) getPath(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewMapClient(cnx)
		rep, err := client.GetPath(ctx, &proto.PathRequest{MapName: args.MapName, Src: args.Src, Dst: args.Dst})
		if err != nil {
			return err
		}

		var out []uint64
		for i := uint32(0); i < args.Max; i++ {
			x, err := rep.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			out = append(out, x.GetId())
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)
		return nil
	})
}

type actionFunc func(ctx context.Context, cli *grpc.ClientConn) error

func (c *ClientCLI) connect(ctx context.Context, action actionFunc) error {
	endpoint, err := discovery.DefaultDiscovery.Map()
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

// PathArgs gathers the possible arguments for the map-related calls.
// Hopefully they are always the same.
type PathArgs struct {
	MapName  string
	Src, Dst uint64
	Max      uint32
}

// Parse extracts the elements of the PathArgs from the array of command line
// positional arguments.
func (pa *PathArgs) Parse(args []string) (err error) {
	if len(args) >= 1 {
		pa.MapName = args[0]
		if len(args) >= 2 {
			pa.Src, err = strconv.ParseUint(args[1], 10, 63)
			if err != nil {
				return err
			}
			if len(args) >= 3 {
				pa.Dst, err = strconv.ParseUint(args[2], 10, 63)
				if err != nil {
					return err
				}
				if len(args) >= 4 {
					return errors.New("Maximum 3 arguments expected: MAPNAME [INT [INT]]")
				}
			}
		}
	} else {
		return errors.New("Minimum 1 argument expected: MAPNAME [INT [INT]]")
	}
	return nil
}
