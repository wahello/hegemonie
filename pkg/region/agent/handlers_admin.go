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

func (srv *srvAdmin) rlockDo(action func() error) error {
	srv.w.RLock()
	defer srv.w.RUnlock()
	return action()
}

func (srv *srvAdmin) wlockDo(action func() error) error {
	srv.w.WLock()
	defer srv.w.WUnlock()
	return action()
}

func (s *srvAdmin) Produce(ctx context.Context, req *proto.None) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(func() error { s.w.Produce(); return nil })
}

func (s *srvAdmin) Move(ctx context.Context, req *proto.None) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(func() error { s.w.Move(); return nil })
}

func (s *srvAdmin) Save(ctx context.Context, req *proto.None) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(func() error { return s.w.SaveLiveToFiles(s.cfg.pathSave) })
}

func (s *srvAdmin) GetScores(ctx context.Context, req *proto.None) (*proto.ListOfCities, error) {
	sb := &proto.ListOfCities{}
	err := s.rlockDo(func() error {
		for _, c := range s.w.Live.Cities {
			sb.Items = append(sb.Items, ShowCityPublic(s.w, c, true))
		}
		return nil
	})
	return sb, err
}
