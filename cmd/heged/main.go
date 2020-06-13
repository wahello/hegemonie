// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	hegemonie_auth_agent "github.com/jfsmig/hegemonie/pkg/auth/agent"
	hegemonie_event_agent "github.com/jfsmig/hegemonie/pkg/event/agent"
	hegemonie_region_agent "github.com/jfsmig/hegemonie/pkg/region/agent"
	"github.com/jfsmig/hegemonie/pkg/utils"
	hegemonie_web_agent "github.com/jfsmig/hegemonie/pkg/web/agent"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	regCmd := hegemonie_region_agent.Command()
	regCmd.Use = "region"
	regCmd.Aliases = []string{"reg"}

	aaaCmd := hegemonie_auth_agent.Command()
	aaaCmd.Use = "auth"
	aaaCmd.Aliases = []string{"aaa"}

	evtCmd := hegemonie_event_agent.Command()
	evtCmd.Use = "evt"
	evtCmd.Aliases = []string{"event", "events"}

	webCmd := hegemonie_web_agent.Command()
	webCmd.Use = "web"
	webCmd.Aliases = []string{"html", "http", "front"}

	rootCmd := &cobra.Command{
		Use:   "heged",
		Short: "Hegemonie main CLI",
		Long:  "Hegemonie: main binary tool to start service agents, query clients and operation jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	rootCmd.AddCommand(aaaCmd)
	rootCmd.AddCommand(regCmd)
	rootCmd.AddCommand(evtCmd)
	rootCmd.AddCommand(webCmd)
	utils.PatchCommandLogs(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
