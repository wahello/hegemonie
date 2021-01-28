// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jfsmig/hegemonie/pkg/auth/client"
	"github.com/jfsmig/hegemonie/pkg/event/agent"
	"github.com/jfsmig/hegemonie/pkg/event/client"
	"github.com/jfsmig/hegemonie/pkg/map/agent"
	"github.com/jfsmig/hegemonie/pkg/map/client"
	"github.com/jfsmig/hegemonie/pkg/region/agent"
	"github.com/jfsmig/hegemonie/pkg/region/client"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type srvCommons struct {
	endpointSrv string
	endpointMon string
	serviceType string
	pathKey     string
	pathCrt     string
}

const (
	defaultKeyPath = "/etc/hegemonie/pki/<SRVTYPE>.key"
	defaultCrtPath = "/etc/hegemonie/pki/<SRVTYPE>.crt"
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

func servers() *cobra.Command {
	var srv srvCommons
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run Hegemonie services",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	cmd.TraverseChildren = true
	cmd.PersistentFlags().StringVar(&srv.pathKey,
		"tls-key", defaultKeyPath, "Path to the X509 key file")
	cmd.PersistentFlags().StringVar(&srv.pathCrt,
		"tls-crt", defaultCrtPath, "Path to the X509 certificate file")
	cmd.PersistentFlags().StringVar(&srv.endpointSrv,
		"endpoint", fmt.Sprintf("0.0.0.0:%v", utils.DefaultPortCommon),
		"IP:PORT endpoint for the gRPC server")
	cmd.PersistentFlags().StringVar(&srv.endpointMon,
		"monitoring", fmt.Sprintf("0.0.0.0:%v", utils.DefaultPortMonitoring),
		"IP:PORT endpoint for the HTTP/1.1 Prometheus exporter")

	cmd.AddCommand(srv.maps(), srv.events(), srv.regions())

	return cmd
}

func clients() *cobra.Command {
	var proxy string

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client tool for various Hegemonie services",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	// Set a common reasonable timeout to all client RPC
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
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
	//cmd.PersistentPostRun = func(_ *cobra.Command, _ []string) { cancel() }

	cmd.AddCommand(clientMap(ctx), clientEvent(ctx), clientAuth(ctx), clientRegion(ctx))
	return cmd
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

func clientMap(ctx context.Context) *cobra.Command {
	var cfg mapclient.ClientCLI
	var pathArgs mapclient.PathArgs

	hook := func(action func() error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			if err := pathArgs.Parse(args); err != nil {
				return errors.Trace(err)
			}
			return action()
		}
	}

	cmd := &cobra.Command{
		Use:   "maps",
		Short: "Client of a Maps API service",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all the maps registered",
		Example: "map list [$MAPID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pathArgs.MapName = first(args)
			return cfg.ListMaps(ctx, pathArgs)
		},
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
		Use:   "events",
		Short: "Client of an Events API service",
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
					return errors.Trace(err)
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
					return errors.Trace(err)
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
		Use:     "regions",
		Short:   "Client of a Regions API service",
		Example: "hege client regions (create|list) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	createRegion := &cobra.Command{
		Use:     "create",
		Short:   "Create a new region",
		Example: "hege client regions create $REGION_ID $MAP_ID",
		Args:    cobra.ExactArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCreateRegion(ctx, args[0], args[1]) },
	}

	listRegions := &cobra.Command{
		Use:     "list",
		Short:   "List the existing regions",
		Example: "hege client regions list [$REGION_ID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoListRegions(ctx, first(args)) },
	}

	roundMovement := &cobra.Command{
		Use:     "move",
		Short:   "Execute a movement round on the region",
		Example: "hege client regions move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionMovement(ctx, args[0]) },
	}

	roundProduction := &cobra.Command{
		Use:     "produce",
		Short:   "Execute a movement round on the region",
		Example: "hege client regions move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionProduction(ctx, args[0]) },
	}

	cmd.AddCommand(createRegion, listRegions, roundMovement, roundProduction)
	return cmd
}

func (srv *srvCommons) events() *cobra.Command {
	return &cobra.Command{
		Use:               "events",
		Short:             "Event Log Service",
		Example:           "hege server events /path/to/event/rocksdb",
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: srv.wrapPreRun("events"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := evtagent.Config{PathBase: args[0]}
			return srv.runServer(cfg)
		},
	}
}

func (srv *srvCommons) maps() *cobra.Command {
	pathMaps := "/etc/hegemonie/maps"
	cmd := &cobra.Command{
		Use:               "maps",
		Short:             "Map Service",
		Example:           "hege server maps",
		Args:              cobra.NoArgs,
		PersistentPreRunE: srv.wrapPreRun("maps"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mapagent.Config{PathRepository: pathMaps}
			return srv.runServer(cfg)
		},
	}
	cmd.PersistentFlags().StringVarP(
		&pathMaps, "defs", "d", pathMaps,
		"Explicit path to the directory with the JSON definitions of the maps")
	return cmd
}

func (srv *srvCommons) regions() *cobra.Command {
	pathDefinitions := "/etc/hegemonie/definitions"
	cmd := &cobra.Command{
		Use:               "regions",
		Short:             "Region Service",
		Example:           "hege server regions /path/to/live/dir",
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: srv.wrapPreRun("regions"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := regagent.Config{PathDefs: pathDefinitions, PathLive: args[0]}
			return srv.runServer(cfg)
		},
	}
	cmd.PersistentFlags().StringVarP(
		&pathDefinitions, "defs", "d", pathDefinitions,
		"Explicit path to the directory with the JSON definitions of the world")
	return cmd
}

func (srv *srvCommons) replaceTag(ps *string) {
	*ps = strings.Replace(*ps, "<SRVTYPE>", srv.serviceType, 1)
}

func nonLeaf(_ *cobra.Command, _ []string) error { return errors.New("missing subcommand") }

func first(args []string) string {
	if len(args) <= 0 {
		return ""
	}
	return args[0]
}

type appRegistrator interface {
	Register(ctx context.Context, grpcServer *grpc.Server) error
}

func (srv *srvCommons) runServer(reg appRegistrator) error {
	var listenerSrv, listenerMon net.Listener
	var grpcSrv *grpc.Server
	var prometheusExporter *http.Server
	var err error

	ctx, cancel := context.WithCancel(context.Background())

	grpcSrv, err = srv.ServerTLS()
	if err != nil {
		return errors.Annotate(err, "TLS server error")
	}

	err = reg.Register(ctx, grpcSrv)
	if err != nil {
		return errors.Annotate(err, "App config error")
	}

	listenerSrv, err = net.Listen("tcp", srv.endpointSrv)
	if err != nil {
		return errors.NewNotValid(err, "listen error")
	}

	if srv.endpointMon != "" {
		listenerMon, err = net.Listen("tcp", srv.endpointMon)
		if err != nil {
			cancel()
			return errors.NewNotValid(err, "listen error")
		}

		prometheusExporter = &http.Server{Handler: promhttp.Handler()}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopChan)

	var barrier sync.WaitGroup
	runner := func(wg *sync.WaitGroup, cb func() error) {
		defer wg.Done()
		if err := cb(); err != http.ErrServerClosed {
			utils.Logger.Error().Err(err).Msg("failed")
		} else {
			utils.Logger.Info().Err(err).Msg("exiting")
		}
		cancel()
	}

	barrier.Add(1)
	go runner(&barrier, func() error { return grpcSrv.Serve(listenerSrv) })

	if prometheusExporter != nil {
		barrier.Add(1)
		go runner(&barrier, func() error { return prometheusExporter.Serve(listenerMon) })
	}

	select {
	case <-stopChan:
		break
	case <-ctx.Done():
		break
	}
	cancel()

	grpcSrv.GracefulStop()

	if prometheusExporter != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := prometheusExporter.Shutdown(shutdownCtx); err != nil {
			utils.Logger.Warn().Err(err).Msg("shutdown error")
		}
	}

	barrier.Wait()
	return nil
}

// ServerTLS automates the creation of a grpc.Server over a TLS connection
// with the proper interceptors.
func (srv *srvCommons) ServerTLS() (*grpc.Server, error) {
	if len(srv.pathCrt) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 certificate path [%s]", srv.pathCrt)
	}
	if len(srv.pathKey) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 key path [%s]", srv.pathKey)
	}
	var certBytes, keyBytes []byte
	var err error

	utils.Logger.Info().Str("key", srv.pathKey).Str("crt", srv.pathCrt).Msg("TLS config")

	if certBytes, err = ioutil.ReadFile(srv.pathCrt); err != nil {
		return nil, errors.Annotate(err, "certificate file error")
	}
	if keyBytes, err = ioutil.ReadFile(srv.pathKey); err != nil {
		return nil, errors.Annotate(err, "key file error")
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, errors.New("invalid certificates")
	}

	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, errors.Annotate(err, "x509 key pair error")
	}

	return grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			utils.NewUnaryServerInterceptorZerolog())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			utils.NewStreamServerInterceptorZerolog()))), nil
}

func (srv *srvCommons) wrapPreRun(srvtype string) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) (err error) {
		srv.serviceType = srvtype
		utils.OverrideLogID("hege," + srv.serviceType)
		utils.ApplyLogModifiers()
		srv.replaceTag(&srv.pathKey)
		srv.replaceTag(&srv.pathCrt)
		return nil
	}
}
