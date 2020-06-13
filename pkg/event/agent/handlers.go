// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_agent

import (
	"context"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"math"
)

func (srv *eventService) Ack1(ctx context.Context, req *proto.Ack1Req) (*proto.None, error) {
	err := srv.backend.Ack1(req.CharId, req.When, req.EvtId)
	return &proto.None{}, err
}

func (srv *eventService) List(ctx context.Context, req *proto.ListReq) (*proto.ListRep, error) {
	items, err := srv.backend.List(req.CharId, req.Marker, req.Max)
	if err != nil {
		return nil, err
	} else {
		rep := proto.ListRep{}
		for _, x := range items {
			rep.Items = append(rep.Items, &proto.ListItem{
				CharId:  x.CharId,
				When:    math.MaxUint64 - x.When,
				EvtId:   x.Id,
				Payload: x.Payload,
			})
		}
		return &rep, nil
	}
}

func (srv *eventService) Push1(ctx context.Context, req *proto.Push1Req) (*proto.None, error) {
	err := srv.backend.Push1(req.CharId, req.EvtId, req.Payload)
	return &proto.None{}, err
}
