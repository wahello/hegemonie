// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	hegemonie_auth_client "github.com/jfsmig/hegemonie/pkg/auth/client"
	hegemonie_event_client "github.com/jfsmig/hegemonie/pkg/event/client"
	"github.com/jfsmig/hegemonie/pkg/region/client"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	regCmd := hegemonie_region_client.Command()
	regCmd.Use = "region"
	regCmd.Aliases = []string{"reg"}

	aaaCmd := hegemonie_auth_client.Command()
	aaaCmd.Use = "auth"
	aaaCmd.Aliases = []string{"aaa"}

	evtCmd := hegemonie_event_client.Command()
	evtCmd.Use = "event"
	evtCmd.Aliases = []string{"evt"}

	rootCmd := &cobra.Command{
		Use:   "hege",
		Short: "Hegemonie CLI",
		Long:  "Hegemonie client with subcommands for the several agents of an Hegemonie system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	rootCmd.AddCommand(aaaCmd, regCmd, evtCmd)
	utils.PatchCommandLogs(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
