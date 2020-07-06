// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
)

type authConfig struct {
	endpoint string
}

func Command() *cobra.Command {
	cfg := authConfig{}

	cmd := &cobra.Command{
		Use:   "region",
		Short: "API service client",
		Args:  cobra.MinimumNArgs(1),
		RunE:  utils.NonLeaf,
	}

	cmd.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.EndpointLocal(utils.DefaultPortRegion), "IP:PORT endpoint for the TCP/IP server")
	return cmd
}
