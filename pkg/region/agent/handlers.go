// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/discovery"
	"github.com/jfsmig/hegemonie/pkg/healthcheck"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc"
	"net"
)

// Config gathers the configuration fields required to start a gRPC region API service.
type Config struct {
	Endpoint string
	Backend  string
}

// Run starts a Region API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (cfg *Config) Run(ctx context.Context) error {
	var err error
	var w region.World

	w.Init()

	if cfg.Backend == "" {
		return errors.New("Missing path for live data")
	}

	err = w.Sections(cfg.Backend).Load()
	if err != nil {
		return err
	}

	err = w.PostLoad()
	if err != nil {
		return fmt.Errorf("Inconsistent World from [%s]: %v", cfg.Backend, err)
	}

	err = w.Check()
	if err != nil {
		return fmt.Errorf("Inconsistent World from [%s]: %v", cfg.Backend, err)
	}

	lis, err := net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var eventEndpoint string
	eventEndpoint, err = discovery.DefaultDiscovery.Event()
	if err != nil {
		return fmt.Errorf("Invalid Event service configured [%s]: %v", eventEndpoint, err)
	}
	var cnxEvent *grpc.ClientConn
	cnxEvent, err = grpc.Dial(eventEndpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cnxEvent.Close()
	w.SetNotifier(&EventStore{cnx: cnxEvent})

	srv := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterCityServer(srv, &srvCity{cfg: cfg, w: &w})
	proto.RegisterDefinitionsServer(srv, &srvDefinitions{cfg: cfg, w: &w})
	proto.RegisterAdminServer(srv, &srvAdmin{cfg: cfg, w: &w})
	proto.RegisterArmyServer(srv, &srvArmy{cfg: cfg, w: &w})
	grpc_health_v1.RegisterHealthServer(srv, &srvHealth{w: &w})

	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
