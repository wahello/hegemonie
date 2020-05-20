// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mapper",
		Short: "Handle map graphs",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}
	rootCmd.AddCommand(CommandNormalize())
	rootCmd.AddCommand(CommandSplit())
	rootCmd.AddCommand(CommandDot())
	rootCmd.AddCommand(CommandSvg())
	rootCmd.AddCommand(CommandExport())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
