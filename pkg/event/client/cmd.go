// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_client

import (
	"errors"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
)

type eventConfig struct {
	endpoint string
}

func Command() *cobra.Command {
	cfg := eventConfig{}

	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cli"},
		Short:   "Event service client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}

	push := &cobra.Command{
		Use:   "push",
		Short: "Push an event",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doPush(cmd, args, &cfg)
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the events",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doList(cmd, args, &cfg)
		},
	}

	ack := &cobra.Command{
		Use:   "ack",
		Short: "Acknowledge an event",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doAck(cmd, args, &cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointEvent, "IP:PORT endpoint for the TCP/IP server")
	cmd.AddCommand(push, ack, list)
	return cmd
}
