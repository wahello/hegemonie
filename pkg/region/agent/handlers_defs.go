// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"io"
)

type srvDefinitions struct {
	cfg *Config
	w   *region.World
}

func (s *srvDefinitions) ListUnits(req *proto.PaginatedQuery, stream proto.Definitions_ListUnitsServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	last := req.GetMarker()
	for {
		tab := s.w.Definitions.Units.Slice(last, 100)
		if len(tab) <= 0 {
			return nil
		}
		for _, i := range tab {
			last = i.ID
			err := stream.Send(&proto.UnitTypeView{
				Id: i.ID, Name: i.Name, Ticks: i.Ticks, Health: i.Health})
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (s *srvDefinitions) ListBuildings(req *proto.PaginatedQuery, stream proto.Definitions_ListBuildingsServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	last := req.GetMarker()
	for {
		tab := s.w.Definitions.Buildings.Slice(last, 100)
		if len(tab) <= 0 {
			return nil
		}
		for _, i := range tab {
			last = i.ID
			err := stream.Send(&proto.BuildingTypeView{
				Id: i.ID, Name: i.Name, Ticks: i.Ticks})
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (s *srvDefinitions) ListKnowledges(req *proto.PaginatedQuery, stream proto.Definitions_ListKnowledgesServer) error {
	s.w.RLock()
	defer s.w.RUnlock()

	last := req.GetMarker()
	for {
		tab := s.w.Definitions.Knowledges.Slice(last, 100)
		if len(tab) <= 0 {
			return nil
		}
		for _, i := range tab {
			last = i.ID
			err := stream.Send(&proto.KnowledgeTypeView{
				Id: i.ID, Name: i.Name, Ticks: i.Ticks})
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}
