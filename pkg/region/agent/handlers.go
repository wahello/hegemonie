// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/healthcheck"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

// Config gathers the configuration fields required to start a gRPC region API service.
type Config struct {
	Endpoint string
	PathDefs string
	PathLive string
}

type regionApp struct {
	cfg *Config
	w   *region.World
}

var none = &proto.None{}

// Run starts a Region API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (cfg *Config) Run(_ context.Context, grpcSrv *grpc.Server) error {
	var err error
	var w region.World

	w.Init()

	if cfg.PathDefs == "" {
		return errors.NotValidf("Missing path for definition data directory")
	}

	if cfg.PathLive == "" {
		return errors.NotValidf("Missing path for live data directory")
	}

	err = w.LoadDefinitions(cfg.PathDefs)
	if err != nil {
		return errors.Annotate(err, "pathDefs error")
	}

	err = w.LoadRegions(cfg.PathLive)
	if err != nil {
		return errors.Annotate(err, "pathLive error")
	}

	err = w.Check()
	if err != nil {
		return errors.Annotate(err, "inconsistent world")
	}

	lis, err := net.Listen("tcp", cfg.Endpoint)
	if err != nil {
		return errors.Annotate(err, "listen error")
	}

	var eventEndpoint string
	eventEndpoint, err = utils.DefaultDiscovery.Event()
	if err != nil {
		return errors.Annotatef(err, "Invalid Event service configured [%s]", eventEndpoint)
	}
	var cnxEvent *grpc.ClientConn
	cnxEvent, err = grpc.Dial(eventEndpoint, grpc.WithInsecure())
	if err != nil {
		return errors.Annotate(err, "dial error")
	}
	defer cnxEvent.Close()
	w.SetNotifier(&EventStore{cnx: cnxEvent})

	grpc_health_v1.RegisterHealthServer(grpcSrv, &regionApp{w: &w, cfg: cfg})
	proto.RegisterCityServer(grpcSrv, &cityApp{regionApp{w: &w, cfg: cfg}})
	proto.RegisterDefinitionsServer(grpcSrv, &defsApp{regionApp{w: &w, cfg: cfg}})
	proto.RegisterAdminServer(grpcSrv, &adminApp{regionApp{w: &w, cfg: cfg}})
	proto.RegisterArmyServer(grpcSrv, &armyApp{regionApp{w: &w, cfg: cfg}})

	utils.Logger.Info().
		Str("defs", cfg.PathDefs).
		Str("live", cfg.PathLive).
		Str("endpoint", cfg.Endpoint).
		Msg("starting")

	return grpcSrv.Serve(lis)
}

func (app *regionApp) _worldLock(mode rune, action func() error) error {
	switch mode {
	case 'r':
		app.w.RLock()
		defer app.w.RUnlock()
	case 'w':
		app.w.WLock()
		defer app.w.WUnlock()
	default:
		return status.Error(codes.Internal, "invalid lock type")
	}
	return action()
}

func (app *regionApp) _regLock(mode rune, regID string, action func(*region.Region) error) error {
	switch mode {
	case 'r':
		app.w.RLock()
		defer app.w.RUnlock()
	case 'w':
		app.w.WLock()
		defer app.w.WUnlock()
	default:
		return status.Error(codes.Internal, "invalid lock type")
	}
	r := app.w.Regions.Get(regID)
	if r == nil {
		return status.Error(codes.NotFound, "no such region")
	}
	return action(r)
}

func (app *regionApp) cityLock(mode rune, req *proto.CityId, action func(*region.Region, *region.City) error) error {
	return app._regLock('r', req.Region, func(r *region.Region) error {
		switch mode {
		case 'r':
			// TODO(jfs) NYI
		case 'w':
			// TODO(jfs) NYI
		default:
			return status.Error(codes.Internal, "invalid lock type")
		}

		c := r.CityGet(req.City)
		if c == nil {
			return status.Error(codes.NotFound, "no such city")
		}
		if c.Deputy != req.Character && c.Owner != req.Character {
			return status.Error(codes.PermissionDenied, "permission denied")
		}

		return action(r, c)
	})
}

func (app *regionApp) armyLock(mode rune, req *proto.ArmyId, action func(*region.Region, *region.City, *region.Army) error) error {
	cID := proto.CityId{Region: req.Region, City: req.City, Character: req.Character}
	return app.cityLock(mode, &cID, func(r *region.Region, c *region.City) error {
		a := c.Armies.Get(req.Army)
		if a == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return action(r, c, a)
	})
}
