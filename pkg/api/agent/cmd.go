// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_api_agent

import (
	"errors"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "agent",
		Aliases: []string{"worker", "server", "service", "srv"},
		Short:   "API service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("NYI")
		},
	}
}
