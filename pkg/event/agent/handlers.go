// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtagent

import (
	"context"
	"errors"
	"fmt"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	"math"
	"net"

	"google.golang.org/grpc"

	back "github.com/jfsmig/hegemonie/pkg/event/backend-local"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
)

// Config gathers the configuration fields required to start a gRPC Event API service.
type Config struct {
	Endpoint string
	PathBase string
}

type eventService struct {
	cfg     Config
	backend *back.Backend
}

// Run starts an Event API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (cfg Config) Run(ctx context.Context) error {
	if cfg.PathBase == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var srv eventService
	var lis net.Listener
	var err error

	srv.cfg = cfg
	srv.backend, err = back.Open(srv.cfg.PathBase)
	if err != nil {
		return err
	}

	lis, err = net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		utils.ServerUnaryInterceptorZerolog(),
		utils.ServerStreamInterceptorZerolog())
	grpc_health_v1.RegisterHealthServer(server, &srv)
	proto.RegisterProducerServer(server, &srv)
	proto.RegisterConsumerServer(server, &srv)

	utils.Logger.Info().
		Str("base", srv.cfg.PathBase).
		Str("url", srv.cfg.Endpoint).
		Msg("starting")
	return server.Serve(lis)
}

// Ack1 marks an event as read so that it won't be listed again.
func (es *eventService) Ack1(ctx context.Context, req *proto.Ack1Req) (*proto.None, error) {
	err := es.backend.Ack1(req.CharId, req.When, req.EvtId)
	return &proto.None{}, err
}

// List streams event objects belonging to the user with the given ID. The objects are sorted by
// decreasing timestamp then by increasing UUID. The events are served as they are stored, the
// messages are not rendered.
func (es *eventService) List(ctx context.Context, req *proto.ListReq) (*proto.ListRep, error) {
	items, err := es.backend.List(req.CharId, req.Marker, req.Max)
	if err != nil {
		return nil, err
	}

	rep := proto.ListRep{}
	for _, x := range items {
		rep.Items = append(rep.Items, &proto.ListItem{
			CharId:  x.CharID,
			When:    math.MaxUint64 - x.When,
			EvtId:   x.ID,
			Payload: x.Payload,
		})
	}
	return &rep, nil
}

// Push1 inserts an event in the log of the Character with the given ID.
// The current timestamp will be used. An UUID will be generated.
func (es *eventService) Push1(ctx context.Context, req *proto.Push1Req) (*proto.None, error) {
	err := es.backend.Push1(req.CharId, req.EvtId, req.Payload)
	return &proto.None{}, err
}

// Check implements the one-shot healthcheck of the gRPC service
func (es *eventService) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// FIXME(jfs): check the service ID
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the long polling healthcheck of the gRPC service
func (es *eventService) Watch(req *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
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
