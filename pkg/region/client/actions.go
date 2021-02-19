// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

// ClientCLI gathers the actions destined to be exposed at the CLI, to manage a region service.
type ClientCLI struct{}

// DoCreateRegion triggers the synchronous creation of a region with the given name, modeled on the named map.
func (cli *ClientCLI) DoCreateRegion(ctx context.Context, regID, mapID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).CreateRegion(ctx, &proto.RegionCreateReq{MapName: mapID, Name: regID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, regID, "created")
	})
}

// DoListRegions dumps to os.Stdout a JSON stream of the known regions, sorted by name
func (cli *ClientCLI) DoListRegions(ctx context.Context, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewAdminClient(cnx).ListRegions(ctx, &proto.RegionListReq{NameMarker: marker})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// DoRegionMovement triggers one round of armies movement on all the cities of the named region
func (cli *ClientCLI) DoRegionMovement(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Move(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Moved")
	})
}

// DoRegionProduction triggers the production of resources on all the cities of the named region
func (cli *ClientCLI) DoRegionProduction(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).Produce(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Produced")
	})
}

// DoRegionPushStats triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionPushStats(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).PushStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Done")
	})
}

type _resourcesAbs struct {
	R0 uint64 `json:"r0"`
	R1 uint64 `json:"r1"`
	R2 uint64 `json:"r2"`
	R3 uint64 `json:"r3"`
	R4 uint64 `json:"r4"`
	R5 uint64 `json:"r5"`
}

type _cityStats struct {
	// Identifier
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// Gauges
	StockCapacity  _resourcesAbs `json:"stockCapacity"`
	StockUsage     _resourcesAbs `json:"stockUsage"`
	ScoreBuilding  uint64        `json:"scoreBuilding"`
	ScoreKnowledge uint64        `json:"scoreKnowledge"`
	ScoreArmy      uint64        `json:"scoreArmy"`
	// Counters
	ResourceProduced _resourcesAbs `json:"resourceProduced"`
	ResourceSent     _resourcesAbs `json:"resourceSent"`
	ResourceReceived _resourcesAbs `json:"resourceReceived"`
	TaxSent          _resourcesAbs `json:"taxSent"`
	TaxReceived      _resourcesAbs `json:"taxReceived"`
	Moves            uint64        `json:"moves"`
	UnitRaised       uint64        `json:"unitRaised"`
	UnitLost         uint64        `json:"unitLost"`
	FightJoined      uint64        `json:"fightJoined"`
	FightLeft        uint64        `json:"fightLeft"`
	FightWon         uint64        `json:"fightWon"`
	FightLost        uint64        `json:"fightLost"`
}

func resProto2Json(in *proto.ResourcesAbs) _resourcesAbs {
	return _resourcesAbs{
		R0: in.R0,
		R1: in.R1,
		R2: in.R2,
		R3: in.R3,
		R4: in.R4,
		R5: in.R5,
	}
}
func statsProto2Json(in *proto.CityStats) _cityStats {
	return _cityStats{
		ID:   in.Id,
		Name: in.Name,

		StockCapacity:  resProto2Json(in.StockCapacity),
		StockUsage:     resProto2Json(in.StockUsage),
		ScoreBuilding:  in.ScoreBuilding,
		ScoreKnowledge: in.ScoreKnowledge,
		ScoreArmy:      in.ScoreArmy,

		ResourceProduced: resProto2Json(in.ResourceProduced),
		ResourceSent:     resProto2Json(in.ResourceSent),
		ResourceReceived: resProto2Json(in.ResourceReceived),
		TaxSent:          resProto2Json(in.TaxSent),
		TaxReceived:      resProto2Json(in.TaxReceived),

		Moves:       in.Moves,
		UnitRaised:  in.UnitRaised,
		UnitLost:    in.UnitLost,
		FightJoined: in.FightJoined,
		FightLeft:   in.FightLeft,
		FightWon:    in.FightWon,
		FightLost:   in.FightLost,
	}
}

// DoRegionGetStats triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionGetStats(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewAdminClient(cnx).GetStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			itf, err := rep.Recv()
			if err != nil {
				return nil, err
			}
			// TODO(jfs): Map the protobuf-generated struct to a fac-simile
			// 			  with "omitempty" flags missing. There are probably
			// 			  better ways to achieve this.
			return statsProto2Json(itf), nil
		})
	})
}

// _publicCity is a variant of proto.PublicCity that doesn't omit empty
// fields. So that the value will me printed on os.Stdout
type _publicCity struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Alignment int32  `json:"alignment"`
	Chaos     int32  `json:"chaos"`
	Politics  uint32 `json:"politics"`
	Cult      uint32 `json:"cult"`
	Ethny     uint32 `json:"ethny"`
	Score     int64  `json:"score"`
}

// DoRegionGetScore triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionGetScores(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewAdminClient(cnx).GetScores(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			pc0, e := out.Recv()
			if e != nil {
				return nil, e
			}
			// TODO(jfs): Map the protobuf-generated struct to a fac-simile
			// 			  with "omitempty" flags missing. There are probably
			// 			  better ways to achieve this.
			return &_publicCity{
				ID:        pc0.Id,
				Name:      pc0.Name,
				Alignment: pc0.Alignment,
				Chaos:     pc0.Chaos,
				Politics:  pc0.Politics,
				Cult:      pc0.Cult,
				Ethny:     pc0.Ethny,
				Score:     pc0.Score,
			}, nil
		})
	})
}

func (cli *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Region()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}
