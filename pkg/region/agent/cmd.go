// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/jfsmig/hegemonie/pkg/region/model"
	admin "github.com/jfsmig/hegemonie/pkg/region/proto_admin"
	army "github.com/jfsmig/hegemonie/pkg/region/proto_army"
	city "github.com/jfsmig/hegemonie/pkg/region/proto_city"
)

type regionConfig struct {
	endpoint string
	pathLoad string
	pathSave string
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
		"endpoint", "127.0.0.1:8080", "IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(&cfg.pathLoad,
		"load", "/data/defs", "File to be loaded")
	agent.Flags().StringVar(&cfg.pathSave,
		"save", "/data/dump", "Directory for persistent")

	return agent
}

func e(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args))
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

	if self.pathLoad != "" {
		type cfgSection struct {
			suffix string
			obj    interface{}
		}
		cfgSections := []cfgSection{
			{"defs.json", &w.Definitions},
			{"map.json", &w.Places},
			{"live.json", &w.Live},
		}
		for _, section := range cfgSections {
			var in io.ReadCloser
			p := self.pathLoad + "/" + section.suffix
			in, err = os.Open(p)
			if err != nil {
				return e("Failed to load the World from [%s]: %s", p, err.Error())
			}
			err = json.NewDecoder(in).Decode(section.obj)
			in.Close()
			if err != nil {
				return e("Failed to load the World from [%s]: %s", p, err.Error())
			}
		}
		err = w.PostLoad()
		if err != nil {
			return e("Inconsistent World from [%s]: %s", self.pathLoad, err.Error())
		}
	}

	err = w.Check()
	if err != nil {
		return e("Inconsistent World: %s", err.Error())
	}

	lis, err := net.Listen("tcp", self.endpoint)
	if err != nil {
		return e("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	city.RegisterCityServer(srv, &srvCity{cfg: self, w: &w})
	admin.RegisterAdminServer(srv, &srvAdmin{cfg: self, w: &w})
	army.RegisterArmyServer(srv, &srvArmy{cfg: self, w: &w})
	if err := srv.Serve(lis); err != nil {
		return e("failed to serve: %v", err)
	}

	if self.pathSave != "" {
		err = self.save(&w)
		if err != nil {
			return e("Failed to save the World at exit: %s", err.Error())
		}
	}

	return nil
}

func (self *regionConfig) save(w *region.World) error {
	if self.pathSave == "" {
		return errors.New("No save path configured")
	}
	p := self.pathSave + "/" + makeSaveFilename()
	p = filepath.Clean(p)
	out, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = w.DumpJSON(out)
	out.Close()
	if err != nil {
		_ = os.Remove(p)
		return err
	}

	latest := self.pathSave + "/latest"
	_ = os.Remove(latest)
	_ = os.Symlink(p, latest)
	return nil
}

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}
