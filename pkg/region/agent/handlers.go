// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"errors"
	"fmt"
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
	PathDefs string
	PathLive string
}

// Run starts a Region API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (cfg *Config) Run(_ context.Context, grpcSrv *grpc.Server) error {
	var err error
	var w region.World

	w.Init()

	if cfg.PathDefs == "" {
		return errors.New("Missing path for definition data directory")
	}

	if cfg.PathLive == "" {
		return errors.New("Missing path for live data directory")
	}

	err = w.LoadDefinitions(cfg.PathDefs)
	if err != nil {
		return err
	}

	err = w.LoadRegions(cfg.PathLive)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var eventEndpoint string
	eventEndpoint, err = utils.DefaultDiscovery.Event()
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

	grpc_health_v1.RegisterHealthServer(grpcSrv, &srvHealth{w: &w})
	proto.RegisterCityServer(grpcSrv, &srvCity{cfg: cfg, w: &w})
	proto.RegisterDefinitionsServer(grpcSrv, &srvDefinitions{cfg: cfg, w: &w})
	proto.RegisterAdminServer(grpcSrv, &srvAdmin{cfg: cfg, w: &w})
	proto.RegisterArmyServer(grpcSrv, &srvArmy{cfg: cfg, w: &w})

	utils.Logger.Info().
		Str("defs", cfg.PathDefs).
		Str("live", cfg.PathLive).
		Str("endpoint", cfg.Endpoint).
		Msg("starting")

	return grpcSrv.Serve(lis)
}
