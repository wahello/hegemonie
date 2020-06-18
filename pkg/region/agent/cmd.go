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

func e(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func (self *regionConfig) execute() error {
	var err error

	w := region.World{}
	w.Init()

	if self.pathSave != "" {
		err = os.MkdirAll(self.pathSave, 0755)
		if err != nil {
			return e("Failed to create [%s]: %s", self.pathSave, err.Error())
		}
	}

	if self.pathLive == "" {
		return e("Missing path for live data")
	}
	if self.pathDefs == "" {
		return e("Missing path for definitions data")
	}

	err = w.LoadDefinitionsFromFiles(self.pathDefs)
	if err != nil {
		return err
	}

	err = w.LoadLiveFromFiles(self.pathLive)
	if err != nil {
		return err
	}

	err = w.PostLoad()
	if err != nil {
		return e("Inconsistent World from [%s] and [%s]: %s", self.pathDefs, self.pathLive, err.Error())
	}

	w.Places.Rehash()

	err = w.Check()
	if err != nil {
		return e("Inconsistent World: %s", err.Error())
	}

	lis, err := net.Listen("tcp", self.endpoint)
	if err != nil {
		return e("failed to listen: %v", err)
	}

	var cnxEvent *grpc.ClientConn
	cnxEvent, err = grpc.Dial(self.endpointEvent, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cnxEvent.Close()
	w.SetNotifier(&EventStore{cnx: cnxEvent})

	srv := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterMapServer(srv, &srvMap{cfg: self, w: &w})
	proto.RegisterCityServer(srv, &srvCity{cfg: self, w: &w})
	proto.RegisterDefinitionsServer(srv, &srvDefinitions{cfg: self, w: &w})
	proto.RegisterAdminServer(srv, &srvAdmin{cfg: self, w: &w})
	proto.RegisterArmyServer(srv, &srvArmy{cfg: self, w: &w})
	if err := srv.Serve(lis); err != nil {
		return e("failed to serve: %v", err)
	}

	if self.pathSave != "" {
		if p, err := w.SaveLiveToFiles(self.pathSave); err != nil {
			return e("Failed to save the World at exit: %s", err.Error())
		} else {
			log.Fatalf("World saved at [%s]", p)
		}
	}

	return nil
}
