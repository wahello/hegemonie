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

type srvAdmin struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvAdmin) Produce(ctx context.Context, req *proto.None) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	s.w.Produce()
	return &proto.None{}, nil
}

func (s *srvAdmin) Move(ctx context.Context, req *proto.None) (*proto.None, error) {
	s.w.WLock()
	defer s.w.WUnlock()

	s.w.Move()
	return &proto.None{}, nil
}

func (s *srvAdmin) GetScores(ctx context.Context, req *proto.None) (*proto.ListOfCities, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	sb := &proto.ListOfCities{}
	for _, c := range s.w.Live.Cities {
		sb.Items = append(sb.Items, ShowCityPublic(s.w, c, true))
	}
	return sb, nil
}

func (s *srvAdmin) Save(ctx context.Context, req *proto.None) (*proto.None, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	if _, err := s.w.SaveLiveToFiles(s.cfg.pathSave); err != nil {
		return nil, err
	}
	return &proto.None{}, nil
}
