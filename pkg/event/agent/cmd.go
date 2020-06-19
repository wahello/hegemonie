// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_agent

import (
	"errors"
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	back "github.com/jfsmig/hegemonie/pkg/event/backend-local"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
)

type eventConfig struct {
	endpoint string
	pathBase string
}

type eventService struct {
	cfg     *eventConfig
	backend *back.Backend
}

func Command() *cobra.Command {
	cfg := eventConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Authentication service",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := eventService{cfg: &cfg}
			return srv.execute()
		},
	}

	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointEvent, "IP:PORT endpoint for the gRPC server")
	agent.Flags().StringVar(&cfg.pathBase,
		"base", "", "Path of the DB")

	return agent
}

func (srv *eventService) execute() error {
	if srv.cfg.pathBase == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var err error
	srv.backend, err = back.Open(srv.cfg.pathBase)
	if err != nil {
		return err
	}

	var lis net.Listener
	if lis, err = net.Listen("tcp", srv.cfg.endpoint); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterProducerServer(server, srv)
	proto.RegisterConsumerServer(server, srv)

	utils.Logger.Warn().
		Str("base", srv.cfg.pathBase).
		Str("url", srv.cfg.endpoint).
		Msg("starting")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
