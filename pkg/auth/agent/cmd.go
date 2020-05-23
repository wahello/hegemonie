// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/auth/model"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
	"os"
)

type authConfig struct {
	endpoint string
	pathLive string
	pathSave string
}

type authService struct {
	proto.AuthServer

	db  auth.Db
	cfg *authConfig
}

func Command() *cobra.Command {
	cfg := authConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Authentication service",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := authService{cfg: &cfg}
			return srv.execute()
		},
	}

	agent.Flags().StringVar(
		&cfg.endpoint, "endpoint", "127.0.0.1:8080",
		"IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(
		&cfg.pathLive, "live", "",
		"Path of the DB backup to load at startup")
	agent.Flags().StringVar(
		&cfg.pathSave, "save", "",
		"Path where to save the DB backup at exit")

	return agent
}

func e(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func (service *authService) execute() error {
	service.db.Init()

	if service.cfg.pathLive == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var p string
	p = service.cfg.pathLive + "/auth.json"
	if in, err := os.Open(p); err != nil {
		return e("Failed to open the DB from [%s]: %s", p, err.Error())
	} else {
		err = json.NewDecoder(in).Decode(&service.db.UsersById)
		_ = in.Close()
		if err != nil {
			return e("Failed to load the DB from [%s]: %s", p, err.Error())
		}
	}

	if err := service.postLoad(); err != nil {
		return e("Inconsistent DB in [%s]: %s", service.cfg.pathLive, err.Error())
	}

	if err := service.db.Check(); err != nil {
		return e("Inconsistent DB: %s", err.Error())
	}

	var lis net.Listener
	var err error
	if lis, err = net.Listen("tcp", service.cfg.endpoint); err != nil {
		return e("failed to listen: %v", err)
	}

	server := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterAuthServer(server, service)
	if err := server.Serve(lis); err != nil {
		return e("failed to serve: %v", err)
	}

	if service.cfg.pathSave != "" {
		if err = service.save(); err != nil {
			return e("Failed to save the DB at exit: %s", err.Error())
		}
	}
	return nil
}

func (srv *authService) postLoad() error {
	return srv.db.ReHash()
}

func (srv *authService) save() error {
	return e("NYI")
}
