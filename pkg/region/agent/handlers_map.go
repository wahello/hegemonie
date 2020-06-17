// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvMap struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvMap) Vertices(ctx context.Context, req *proto.PaginatedQuery) (*proto.ListOfVertices, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	rep := &proto.ListOfVertices{}
	for _, x := range s.w.Places.Cells.Slice(req.Marker, req.Max) {
		rep.Items = append(rep.Items, &proto.Vertex{
			Id: x.Id, X: x.X, Y: x.Y, CityId: x.City})
	}
	return rep, nil
}

func (s *srvMap) Edges(ctx context.Context, req *proto.ListEdgesReq) (*proto.ListOfEdges, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	rep := &proto.ListOfEdges{}
	for _, x := range s.w.Places.Roads.Slice(req.MarkerSrc, req.MarkerDst, req.Max) {
		rep.Items = append(rep.Items, &proto.Edge{Src: x.S, Dst: x.D})
	}
	return rep, nil
}

func (s *srvMap) Cities(ctx context.Context, req *proto.PaginatedQuery) (*proto.ListOfCities, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	rep := &proto.ListOfCities{}
	for _, x := range s.w.Live.Cities.Slice(req.Marker, req.Max) {
		rep.Items = append(rep.Items, ShowCityPublic(s.w, x, false))
	}
	return rep, nil
}
