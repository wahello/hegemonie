// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"google.golang.org/grpc"
)

// ClientCLI gathers the event-related client actions available at the command line.
type ClientCLI struct{}

func (cfg *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Event()
	if err != nil {
		return err
	}
	return utils.Connect(ctx, endpoint, action)
}

// DoPush insert an event whose content and target are described on the command line.
// A descriptive error is returned in case of failure.
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoPush(ctx context.Context, charID string, msg ...string) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var anyError bool
		client := proto.NewProducerClient(cnx)
		for _, a := range msg {
			id := uuid.New().String()
			_, err := client.Push1(ctx, &proto.Push1Req{CharId: charID, EvtId: id, Payload: []byte(a)})
			if err != nil {
				anyError = true
				utils.Logger.Error().Str("char", charID).Str("msg", a).Str("uuid", id).Err(err).Msg("PUSH")
			} else {
				utils.Logger.Info().Str("char", charID).Str("msg", a).Str("uuid", id).Msg("PUSH")
			}
		}
		if !anyError {
			return nil
		}
		return errors.New("Errors occured")
	})
}

// DoAck consumes a message whose owner, timestamp and ID are described on the command line.
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoAck(ctx context.Context, charID, evtID string, when uint64) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewConsumerClient(cnx)
		_, err := client.Ack1(ctx, &proto.Ack1Req{CharId: charID, When: when, EvtId: evtID})
		if err != nil {
			return err
		}
		utils.Logger.Info().
			Str("char", charID).
			Uint64("when", when).
			Str("uuid", evtID).
			Msg("ACK")
		return nil
	})
}

// DoList dumps to os.Stdout the Event objects streamed by the contacted service. The output consists
// in a JSON stream of objects separated by a CRLF (i.e. one object per line)
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoList(ctx context.Context, charID string, when uint64, marker string, max uint32) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewConsumerClient(cnx)
		rep, err := client.List(ctx, &proto.ListReq{CharId: charID, Marker: when, Max: 100})
		if err != nil {
			return err
		}
		anyError := false
		for _, x := range rep.Items {
			fmt.Printf("%s %d %s %s\n", x.CharId, x.When, x.EvtId, x.Payload)
		}
		if anyError {
			return errors.New("Invalid events matched")
		}
		return nil
	})
}
