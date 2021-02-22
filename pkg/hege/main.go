// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/map/client"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	cmd := &cobra.Command{
		Use:   "hege",
		Short: "Hegemonie CLI",
		Long:  "Hegemonie client with subcommands for the several agents of an Hegemonie system.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	utils.AddLogFlagsToCommand(cmd)
	cmd.AddCommand(clients(), servers(), tools())
	if err := cmd.Execute(); err != nil {
		log.Fatalln(errors.ErrorStack(err))
	}
}

func tools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Miscellaneous tools to help the operations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	ctx := context.Background()
	cmd.AddCommand(toolsMap(ctx))
	return cmd
}

func toolsMap(_ context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map",
		Short: "Map handling tools",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	normalize := &cobra.Command{
		Use:   "normalize",
		Short: "Normalize the positions in a map (stdin/stdout)",
		Long:  `Read the map description on the standard input, remap the positions of the vertices in the map graph so that they fit in the given boundaries and dump it to the standard output.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return mapclient.ToolNormalize() },
	}

	var maxDist float64
	split := &cobra.Command{
		Use:   "split",
		Short: "Split the long edges of a map (stdin/stdout)",
		Long:  `Read the map on the standard input, split all the edges that are longer to the given value and dump the new graph on the standard output.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return mapclient.ToolSplit(maxDist) },
	}
	split.Flags().Float64VarP(&maxDist, "dist", "d", 60, "Max road length")

	var noise float64
	noisify := &cobra.Command{
		Use:   "noise",
		Short: "Apply a noise on the positon of the nodes (stdin/stdout)",
		Long:  `Read the map on the standard input, randomly alter the positions of the nodes and dump the new graph on the standard output.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return mapclient.ToolNoise(noise) },
	}
	noisify.Flags().Float64VarP(&noise, "noise", "n", 15, "Percent of the image dimension used as max noise variation on non-city nodes positions")

	drawDot := &cobra.Command{
		Use:   "dot",
		Short: "Convert the JSON map to DOT (stdin/stdout)",
		RunE:  func(cmd *cobra.Command, args []string) error { return mapclient.ToolDot() },
	}

	drawSvg := &cobra.Command{
		Use:   "svg",
		Short: "Convert the JSON map to SVG  (stdin/stdout)",
		RunE:  func(cmd *cobra.Command, args []string) error { return mapclient.ToolSvg() },
	}

	seedInit := &cobra.Command{
		Use:     "init",
		Aliases: []string{"seed"},
		Short:   "Convert the JSON map seed to a JSON raw map (stdin/stdout)",
		RunE:    func(cmd *cobra.Command, args []string) error { return mapclient.ToolInit() },
	}

	cmd.AddCommand(normalize, split, noisify, drawDot, drawSvg, seedInit)
	return cmd
}

func nonLeaf(_ *cobra.Command, _ []string) error { return errors.New("missing subcommand") }
