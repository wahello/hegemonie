// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_client

import (
	"errors"
	"github.com/spf13/cobra"
)

type authConfig struct {
	endpoint string
}

func Command() *cobra.Command {
	cfg := authConfig{}

	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cli"},
		Short:   "Auth service client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}

	show := &cobra.Command{
		Use:     "show",
		Aliases: []string{"check", "view"},
		Short:   "Show the details of a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doShow(cmd, args, &cfg)
		},
	}

	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "put"},
		Short:   "Create a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doCreate(cmd, args, &cfg)
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "all"},
		Short:   "List the existung users",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doList(cmd, args, &cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.endpoint, "endpoint", "127.0.0.1:8082", "IP:PORT endpoint for the TCP/IP server")
	cmd.AddCommand(create, show, list)
	return cmd
}
