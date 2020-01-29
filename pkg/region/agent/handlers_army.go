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

	proto "github.com/jfsmig/hegemonie/pkg/region/proto_army"
)

type srvArmy struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvArmy) Show(ctx context.Context, req *proto.ArmyId) (*proto.ArmyView, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvArmy) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvArmy) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}

func (s *srvArmy) Command(ctx context.Context, req *proto.ArmyCommandReq) (*proto.None, error) {
	return nil, status.Errorf(codes.Unimplemented, "NYI")
}
