// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/healthcheck"
	mapgraph "github.com/jfsmig/hegemonie/pkg/map/graph"
	"github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Config gathers the configuration fields required to start a gRPC map API service.
type Config struct {
	Endpoint       string
	PathRepository string
}

type srvMap struct {
	config *Config
	maps   mapgraph.SetOfMaps
	rw     sync.RWMutex
}

// Run starts an Map API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (cfg *Config) Run(_ context.Context) error {
	srv := &srvMap{config: cfg, maps: make(mapgraph.SetOfMaps, 0)}
	if err := srv.LoadDirectory(cfg.PathRepository); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		utils.ServerUnaryInterceptorZerolog(),
		utils.ServerStreamInterceptorZerolog())
	grpc_health_v1.RegisterHealthServer(grpcServer, srv)
	proto.RegisterMapServer(grpcServer, srv)

	utils.Logger.Info().
		Int("maps", srv.maps.Len()).
		Str("endpoint", cfg.Endpoint).
		Msg("Starting")
	for _, m := range srv.maps {
		utils.Logger.Debug().
			Str("name", m.ID).
			Int("sites", m.Cells.Len()).
			Int("roads", m.Roads.Len()).
			Msg("map>")
	}
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Check implements the one-shot healthcheck of the gRPC service
func (s *srvMap) Check(_ context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// FIXME(jfs): check the service ID
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the long polling healthcheck of the gRPC service
func (s *srvMap) Watch(_ *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	// FIXME(jfs): check the service ID
	for {
		err := srv.Send(&grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		})
		if err != nil {
			return err
		}
	}
}

// Vertices streams Vertice objects, sorted by ID.
func (s *srvMap) Vertices(req *proto.ListVerticesReq, stream proto.Map_VerticesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("no such map")
	}

	next := req.Marker
	for {
		vertices := m.Cells.Slice(next, 100)
		if len(vertices) <= 0 {
			return nil
		}
		for _, x := range vertices {
			err := stream.Send(&proto.Vertex{Id: x.ID, X: x.X, Y: x.Y})
			if err != nil {
				return err
			}
			next = x.ID
		}
	}
}

// Edges streams Edge objects, sorted by source then by destination.
func (s *srvMap) Edges(req *proto.ListEdgesReq, stream proto.Map_EdgesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("no such map")
	}

	src, dst := req.MarkerSrc, req.MarkerDst
	for {
		edges := m.Roads.Slice(src, dst, 100)
		if len(edges) <= 0 {
			return nil
		}
		for _, x := range edges {
			err := stream.Send(&proto.Edge{Src: x.S, Dst: x.D})
			if err != nil {
				return err
			}
			src, dst = x.S, x.D
		}
	}
}

// GetPath streams the Vertice elements of the path from the source to the destination.
func (s *srvMap) GetPath(req *proto.PathRequest, stream proto.Map_GetPathServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("no such map")
	}

	src := req.Src
	for {
		next, err := m.PathNextStep(src, req.Dst)
		if err != nil {
			return err
		}
		err = stream.Send(&proto.PathElement{Id: src})
		if err != nil {
			return err
		}
		if next == req.Dst {
			return nil
		}
		src = next
	}
}

// Cities streams City <ID,name> pair objects
func (s *srvMap) Cities(req *proto.ListCitiesReq, stream proto.Map_CitiesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("no such map")
	}

	next := req.Marker
	for {
		cities := m.Cells.Slice(next, 100)
		if len(cities) <= 0 {
			return nil
		}
		for _, v := range cities {
			if v.City != "" {
				err := stream.Send(&proto.CityLocation{
					Id: v.ID, Name: v.City,
				})
				if err != nil {
					return err
				}
			}
			next = v.ID
		}
	}
}

// Maps streams the name of the maps registered in the current service
func (s *srvMap) Maps(req *proto.ListMapsReq, stream proto.Map_MapsServer) error {
	slice := func(marker string) []proto.MapName {
		s.rw.RLock()
		defer s.rw.RUnlock()
		out := make([]proto.MapName, 0)
		for _, m := range s.maps.Slice(marker, 100) {
			out = append(out, proto.MapName{
				Name:          m.ID,
				CountEdges:    uint32(len(m.Roads)),
				CountVertices: uint32(len(m.Cells)),
				CountCities: func() (total uint32) {
					for _, c := range m.Cells {
						if c.City != "" {
							total++
						}
					}
					return total
				}(),
			})
		}
		return out
	}

	next := req.Marker
	for {
		names := slice(next)
		if len(names) <= 0 {
			return nil
		}
		for _, v := range names {
			if err := stream.Send(&v); err != nil {
				return err
			}
			next = v.Name
		}
	}
}

// LoadDirectory loads all the maps stored as files, containing JSON objects desribing maps.
// Only the filenames with a .final.json suffix are considered.
func (s *srvMap) LoadDirectory(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only accept non-hidden JSON files
		_, fn := filepath.Split(path)
		if info.IsDir() || info.Size() <= 0 {
			return nil
		}
		if len(fn) < 2 || fn[0] == '.' {
			return nil
		}
		if !strings.HasSuffix(fn, ".final.json") {
			return nil
		}

		m := mapgraph.NewMap()
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = m.Load(f); err != nil {
			return err
		}

		s.maps.Add(m)
		return nil
	})
}
