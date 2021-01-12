// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/jfsmig/hegemonie/pkg/auth/client"
	"github.com/jfsmig/hegemonie/pkg/event/agent"
	"github.com/jfsmig/hegemonie/pkg/event/client"
	"github.com/jfsmig/hegemonie/pkg/map/agent"
	"github.com/jfsmig/hegemonie/pkg/map/client"
	"github.com/jfsmig/hegemonie/pkg/region/agent"
	"github.com/jfsmig/hegemonie/pkg/region/client"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	cmd := &cobra.Command{
		Use:   "hege",
		Short: "Hegemonie CLI",
		Long:  "Hegemonie client with subcommands for the several agents of an Hegemonie system.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	utils.PatchCommandLogs(cmd)
	ctx := context.Background()
	cmd.AddCommand(clients(ctx), servers(ctx), tools(ctx))
	if err := cmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}

type srvCommons struct {
	pathKey string
	pathCrt string
	grpcSrv *grpc.Server
}

func servers(ctx context.Context) *cobra.Command {
	var srv srvCommons
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run Hegemonie services",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	cmd.PersistentFlags().StringVar(&srv.pathKey, "key", "", "Path to the X509 key file")
	cmd.PersistentFlags().StringVar(&srv.pathCrt, "crt", "", "Path to the X509 cert file")
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		srv.grpcSrv, err = utils.ServerTLS(srv.pathKey, srv.pathCrt)
		return err
	}
	cmd.AddCommand(srv.maps(ctx), srv.event(ctx), srv.region(ctx))
	return cmd
}

func clients(ctx context.Context) *cobra.Command {
	var proxy string

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client tool for various Hegemonie services",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	// Set a common reasonable timeout to all client RPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	sessionID := os.Getenv("HEGE_CLI_SESSIONID")
	if sessionID == "" {
		sessionID = "cli/" + uuid.New().String()
	}

	// Inherit a session-id from the env
	ctx = metadata.AppendToOutgoingContext(ctx, "session-id", sessionID)

	// Override the discovery if a proxy is configured
	cmd.PersistentFlags().StringVar(&proxy,
		"proxy", "", "IP:PORT endpoint for the gRPC proxy")
	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		if proxy != "" {
			utils.DefaultDiscovery = utils.SingleEndpoint(proxy)
		}
		return nil
	}
	cmd.PersistentPostRun = func(_ *cobra.Command, _ []string) { cancel() }

	cmd.AddCommand(clientMap(ctx), clientEvent(ctx), clientAuth(ctx), clientRegion(ctx))
	return cmd
}

func tools(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Miscellaneous tools to help the operations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
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

func clientMap(ctx context.Context) *cobra.Command {
	var cfg mapclient.ClientCLI
	var pathArgs mapclient.PathArgs

	hook := func(action func() error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			if err := pathArgs.Parse(args); err != nil {
				return err
			}
			return action()
		}
	}

	cmd := &cobra.Command{
		Use:   "map",
		Short: "Map service client",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all the maps registered",
		Example: "map list [$MAPID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE:    hook(func() error { return cfg.ListMaps(ctx, pathArgs) }),
	}

	path := &cobra.Command{
		Use:     "path",
		Short:   "Compute the path between two nodes",
		Example: "map path $MAPID $SRC $DST",
		Args:    cobra.ExactArgs(3),
		RunE:    hook(func() error { return cfg.GetPath(ctx, pathArgs) }),
	}
	path.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "Max path length")

	step := &cobra.Command{
		Use:     "step",
		Short:   "Get the next step of the path between two nodes",
		Example: "map step $REGION $SRC $DST",
		Args:    cobra.ExactArgs(3),
		RunE:    hook(func() error { return cfg.GetStep(ctx, pathArgs) }),
	}

	cities := &cobra.Command{
		Use:     "cities",
		Short:   "List the Cities when the map is instantiated",
		Example: "map cities $REGION [$MARKER]",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    hook(func() error { return cfg.GetCities(ctx, pathArgs) }),
	}
	cities.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N cities")

	roads := &cobra.Command{
		Use:     "roads",
		Short:   "List of the roads of the map",
		Example: "map roads $REGION [$MARKER_SRC [$MARKER_DST]]",
		Args:    cobra.RangeArgs(1, 3),
		RunE:    hook(func() error { return cfg.GetRoads(ctx, pathArgs) }),
	}
	roads.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N roads")

	positions := &cobra.Command{
		Use:     "positions",
		Short:   "List the positions of the map",
		Example: "map positions $REGION [$MARKER]",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    hook(func() error { return cfg.GetPositions(ctx, pathArgs) }),
	}
	positions.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N positions")

	cmd.AddCommand(list, path, step, cities, roads, positions)
	return cmd
}

func clientEvent(ctx context.Context) *cobra.Command {
	var max uint32
	var cfg evtclient.ClientCLI

	cmd := &cobra.Command{
		Use:   "event",
		Short: "Event service client",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	push := &cobra.Command{
		Use:     "push",
		Short:   "Push events in the Character's log",
		Example: `server event push "${CHARACTER}" "${MSG0}" "${MSG1}"`,
		Args:    cobra.MinimumNArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoPush(ctx, args[0], args[1:]...) },
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List the events",
		Example: `server event list "${CHARACTER}" "${EVENT_TIMESTAMP}" [${EVENT_MARKER}]`,
		Args:    cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var when uint64
			var marker string
			if len(args) > 1 {
				var err error
				when, err = strconv.ParseUint(args[1], 10, 63)
				if err != nil {
					return err
				}
				if len(args) > 2 {
					marker = args[2]
				}
			}
			return cfg.DoList(ctx, args[0], when, marker, max)
		},
	}
	list.Flags().Uint32VarP(&max, "max", "m", 0, "List at most N events")

	ack := &cobra.Command{
		Use:     "ack",
		Short:   "Acknowledge an event",
		Example: `server event ack "${CHARACTER}" "${EVENT_UUID}" "${EVENT_TIMESTAMP}"`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var when uint64
			if len(args) > 2 {
				var err error
				when, err = strconv.ParseUint(args[1], 10, 63)
				if err != nil {
					return err
				}
			}
			return cfg.DoAck(ctx, args[0], args[1], when)
		},
	}

	cmd.AddCommand(push, ack, list)
	return cmd

}

func clientAuth(ctx context.Context) *cobra.Command {
	var cfg authclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authorization and Authentication client",
		Example: "auth (users|details|create|invite|affect) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	users := &cobra.Command{
		Use:     "users",
		Short:   "List the registered USERS",
		Example: "auth list",
		Args:    cobra.NoArgs,
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoList(ctx, args) },
	}

	details := &cobra.Command{
		Use:     "detail",
		Short:   "Show the details of specific users",
		Long:    "Print a detailed JSON representation of the information and permissions for each user specified as a positional argument",
		Example: "show a4ddeee6-b72a-4a27-8e2d-35c3cc62c7d3 ab2bca77-efdb-4dc2-b80a-fc03e0fc5226 ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoShow(ctx, args) },
	}

	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a User",
		Example: "auth create forced.user@example.com",
		Args:    cobra.MinimumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCreate(ctx, args) },
	}

	invite := &cobra.Command{
		Use:     "invite",
		Short:   "Invite a user identified by its email",
		Example: "auth invite invited.user@example.com",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoInvite(ctx, args) },
	}

	affect := &cobra.Command{
		Use:   "affect",
		Short: "Invite a user identified by its email",
		RunE:  func(cmd *cobra.Command, args []string) error { return cfg.DoInvite(ctx, args) },
	}

	cmd.AddCommand(users, details, create, invite, affect)
	return cmd
}

func clientRegion(ctx context.Context) *cobra.Command {
	var cfg regclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "region",
		Short:   "Region API client",
		Example: "region (create|list) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	createRegion := &cobra.Command{
		Use:     "create",
		Short:   "Create a new region",
		Example: "region create $REGION_ID $MAP_ID",
		Args:    cobra.ExactArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCreateRegion(ctx, args[0], args[1]) },
	}

	listRegions := &cobra.Command{
		Use:     "list",
		Short:   "List the existing regions",
		Example: "region list [$REGION_ID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoListRegions(ctx, first(args)) },
	}

	roundMovement := &cobra.Command{
		Use:     "move",
		Short:   "Execute a movement round on the region",
		Example: "region move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionMovement(ctx, args[0]) },
	}

	roundProduction := &cobra.Command{
		Use:     "produce",
		Short:   "Execute a movement round on the region",
		Example: "region move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionProduction(ctx, args[0]) },
	}

	cmd.AddCommand(createRegion, listRegions, roundMovement, roundProduction)
	return cmd
}

func (srv *srvCommons) event(ctx context.Context) *cobra.Command {
	var cfg evtagent.Config

	agent := &cobra.Command{
		Use:     "event",
		Short:   "Event Log Service",
		Example: "heged event --endpoint=10.0.0.1:2345 /path/to/event/rocksdb",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.PathBase = args[0]
			return cfg.Run(ctx, srv.grpcSrv)
		},
	}

	agent.Flags().StringVar(&cfg.Endpoint,
		"endpoint", utils.EndpointLocal(utils.DefaultPortEvent), "IP:PORT endpoint for the gRPC server")
	return agent
}

func (srv *srvCommons) maps(ctx context.Context) *cobra.Command {
	var cfg mapagent.Config

	agent := &cobra.Command{
		Use:     "map",
		Short:   "Map Service",
		Example: "heged map --endpoint=10.0.0.1:1234 /path/to/maps/directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.PathRepository = args[0]
			return cfg.Run(ctx, srv.grpcSrv)
		},
	}
	agent.Flags().StringVar(&cfg.Endpoint,
		"endpoint", utils.EndpointLocal(utils.DefaultPortMap), "IP:PORT endpoint for the gRPC server")

	return agent
}

func (srv *srvCommons) region(ctx context.Context) *cobra.Command {
	var cfg regagent.Config

	agent := &cobra.Command{
		Use:     "region",
		Short:   "Region Service",
		Example: "heged region --Endpoint=10.0.0.1:1234 /path/to/defs/dir /path/to/live/dir",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.PathDefs = args[0]
			cfg.PathLive = args[1]
			return cfg.Run(ctx, srv.grpcSrv)
		},
	}
	agent.Flags().StringVar(&cfg.Endpoint,
		"endpoint", utils.EndpointLocal(utils.DefaultPortMap), "IP:PORT Endpoint for the gRPC server")

	return agent
}

func nonLeaf(_ *cobra.Command, _ []string) error {
	return errors.New("missing subcommand")
}
func first(args []string) string {
	if len(args) <= 0 {
		return ""
	}
	return args[0]
}
