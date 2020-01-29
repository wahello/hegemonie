// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_api_client

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
		Short:   "API service client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}

	cmd.Flags().StringVar(&cfg.endpoint, "endpoint", "127.0.0.1:8080", "IP:PORT endpoint for the TCP/IP server")
	return cmd
}
