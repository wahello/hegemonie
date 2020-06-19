// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"errors"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type regionConfig struct {
	endpoint      string
	endpointEvent string
	pathSave      string
	pathDefs      string
	pathLive      string
}

func Command() *cobra.Command {
	cfg := regionConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srvCity", "srv", "service"},
		Short:   "Region service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.execute()
		},
	}
	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointRegion, "IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(&cfg.endpointEvent,
		"event", utils.DefaultEndpointEvent, "Address of the Event server to connect to.")
	agent.Flags().StringVar(&cfg.pathSave,
		"save", "", "Path of the directory where persist the dump of the region.")
	agent.Flags().StringVar(&cfg.pathDefs,
		"defs", "", "Path to the file with the definition of the world.")
	agent.Flags().StringVar(&cfg.pathLive,
		"live", "", "Path to the file with the state of the region.")

	return agent
}

func (cfg *regionConfig) execute() error {
	var err error

	w := region.World{}
	w.Init()

	if cfg.pathSave != "" {
		err = os.MkdirAll(cfg.pathSave, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create [%s]: %v", cfg.pathSave, err)
		}
	}

	if cfg.pathLive == "" {
		return errors.New("Missing path for live data")
	}
	if cfg.pathDefs == "" {
		return errors.New("Missing path for definitions data")
	}

	err = w.LoadDefinitionsFromFiles(cfg.pathDefs)
	if err != nil {
		return err
	}

	err = w.LoadLiveFromFiles(cfg.pathLive)
	if err != nil {
		return err
	}

	err = w.PostLoad()
	if err != nil {
		return fmt.Errorf("Inconsistent World from [%s] and [%s]: %v", cfg.pathDefs, cfg.pathLive, err)
	}

	w.Places.Rehash()

	err = w.Check()
	if err != nil {
		return fmt.Errorf("Inconsistent World: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var cnxEvent *grpc.ClientConn
	cnxEvent, err = grpc.Dial(cfg.endpointEvent, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cnxEvent.Close()
	w.SetNotifier(&EventStore{cnx: cnxEvent})

	srv := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterMapServer(srv, &srvMap{cfg: cfg, w: &w})
	proto.RegisterCityServer(srv, &srvCity{cfg: cfg, w: &w})
	proto.RegisterDefinitionsServer(srv, &srvDefinitions{cfg: cfg, w: &w})
	proto.RegisterAdminServer(srv, &srvAdmin{cfg: cfg, w: &w})
	proto.RegisterArmyServer(srv, &srvArmy{cfg: cfg, w: &w})
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	if cfg.pathSave != "" {
		p, err := w.SaveLiveToFiles(cfg.pathSave)
		if err != nil {
			return fmt.Errorf("Failed to save the World at exit: %v", err)
		}
		log.Fatalf("World saved at [%s]", p)
	}

	return nil
}
