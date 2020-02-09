// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

type srvAdmin struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvAdmin) Produce(ctx context.Context, req *proto.None) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvAdmin) Move(ctx context.Context, req *proto.None) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvAdmin) GetScores(ctx context.Context, req *proto.None) (*proto.ScoreBoard, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvAdmin) Save(ctx context.Context, req *proto.None) (*proto.None, error) {
	s.w.RLock()
	defer s.w.RUnlock()

	err := s.cfg.save(s.w)
	if err != nil {
		return nil, err
	} else {
		return &proto.None{}, nil
	}
}
