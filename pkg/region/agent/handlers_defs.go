// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto_city"
)

type srvDefinitions struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvDefinitions) ListUnits(ctx context.Context, req *proto.PaginatedQuery) (*proto.ListOfUnitTypes, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	v := s.w.Definitions.Units.Slice(req.Marker, ClampU32(req.Max, 1, 1000))
	rep := &proto.ListOfUnitTypes{}
	for _, i := range v {
		rep.Items = append(rep.Items, &proto.UnitTypeView{Id: i.Id, Name: i.Name})
	}
	return rep, nil
}

func (s *srvDefinitions) ListBuildings(ctx context.Context, req *proto.PaginatedQuery) (*proto.ListOfBuildingTypes, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	v := s.w.Definitions.Buildings.Slice(req.Marker, ClampU32(req.Max, 1, 1000))
	rep := &proto.ListOfBuildingTypes{}
	for _, i := range v {
		rep.Items = append(rep.Items, &proto.BuildingTypeView{Id: i.Id, Name: i.Name})
	}
	return rep, nil
}

func (s *srvDefinitions) ListKnowledges(ctx context.Context, req *proto.PaginatedQuery) (*proto.ListOfKnowledgeTypes, error) {
	latch := s.w.ReadLocker()
	latch.Lock()
	defer latch.Unlock()

	v := s.w.Definitions.Knowledges.Slice(req.Marker, ClampU32(req.Max, 1, 1000))
	rep := &proto.ListOfKnowledgeTypes{}
	for _, i := range v {
		rep.Items = append(rep.Items, &proto.KnowledgeTypeView{Id: i.Id, Name: i.Name})
	}
	return rep, nil
}

func ClampU32(v, min, max uint32) uint32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
