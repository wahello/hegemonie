// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"github.com/jfsmig/hegemonie/pkg/api/agent"
	"github.com/jfsmig/hegemonie/pkg/api/client"
	"github.com/jfsmig/hegemonie/pkg/auth/agent"
	"github.com/jfsmig/hegemonie/pkg/auth/client"
	"github.com/jfsmig/hegemonie/pkg/events/agent"
	"github.com/jfsmig/hegemonie/pkg/events/client"
	"github.com/jfsmig/hegemonie/pkg/region/agent"
	"github.com/jfsmig/hegemonie/pkg/region/client"
	"github.com/jfsmig/hegemonie/pkg/web/agent"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	regCmd := &cobra.Command{
		Use:     "region",
		Aliases: []string{"reg"},
		Short:   "Region group of commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	regCmd.AddCommand(hegemonie_region_agent.Command())
	regCmd.AddCommand(hegemonie_region_client.Command())

	aaaCmd := &cobra.Command{
		Use:     "auth",
		Aliases: []string{"aaa"},
		Short:   "Auth group of commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	aaaCmd.AddCommand(hegemonie_auth_agent.Command())
	aaaCmd.AddCommand(hegemonie_auth_client.Command())

	webCmd := &cobra.Command{
		Use:     "web",
		Aliases: []string{"html", "http", "front"},
		Short:   "Web group of commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	webCmd.AddCommand(hegemonie_web_agent.Command())

	rootCmd := &cobra.Command{
		Use:   "hegemonie",
		Short: "Hegemonie main CLI",
		Long:  "Hegemonie: main binary tool to start service agents, query clients and operation jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	rootCmd.AddCommand(aaaCmd)
	rootCmd.AddCommand(regCmd)
	rootCmd.AddCommand(webCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
