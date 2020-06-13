// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_client

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
)

func doPush(cmd *cobra.Command, args []string, cfg *eventConfig) error {
	if len(args) < 2 {
		return errors.New("Too few arguments (minimum: 2)")
	}

	charId, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	ctx := context.Background()

	cnx, err := grpc.Dial(cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewProducerClient(cnx)

	anyError := false
	for _, a := range args[1:] {
		id := uuid.New().String()
		req := proto.Push1Req{
			CharId:  charId,
			EvtId:   id,
			Payload: []byte(a),
		}
		_, err := client.Push1(ctx, &req)
		if err != nil {
			anyError = true
			utils.Logger.Error().Uint64("char", charId).Str("msg", a).Str("uuid", id).Err(err).Msg("PUSH")
		} else {
			utils.Logger.Info().Uint64("char", charId).Str("msg", a).Str("uuid", id).Msg("PUSH")
		}
	}
	if !anyError {
		return nil
	}
	return errors.New("Errors occured")
}

func doAck(cmd *cobra.Command, args []string, cfg *eventConfig) error {
	var charId, when uint64
	var id string
	var err error

	// Parse the input
	if len(args) != 3 {
		return errors.New("3 arguments expected (Character When Uuid)")
	}
	charId, err = strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}
	when, err = strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return err
	}
	id = args[2]

	// Send the request
	ctx := context.Background()
	cnx, err := grpc.Dial(cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewConsumerClient(cnx)
	req := proto.Ack1Req{CharId: charId, When: when, EvtId: id}
	_, err = client.Ack1(ctx, &req)

	if err != nil {
		return err
	}
	utils.Logger.Info().Uint64("char", charId).Uint64("when", when).Str("uuid", id).Msg("ACK")
	return nil
}

func doList(cmd *cobra.Command, args []string, cfg *eventConfig) error {
	var charId, when uint64
	var err error

	// Parse the input
	switch len(args) {
	case 1:
		charId, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return err
		}
	case 2:
		charId, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return err
		}
		when, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.New("1 or 2 arguments expected (CHARACTER [MARKER])")
	}

	// Send the request
	ctx := context.Background()
	cnx, err := grpc.Dial(cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewConsumerClient(cnx)
	req := proto.ListReq{CharId: charId, Marker: when, Max: 100}
	rep, err := client.List(ctx, &req)

	if err != nil {
		return err
	}
	anyError := false
	for _, x := range rep.Items {
		fmt.Printf("%d %d %s %s\n", x.CharId, x.When, x.EvtId, x.Payload)
	}
	if anyError {
		return errors.New("Invalid events matched")
	}
	return nil
}
