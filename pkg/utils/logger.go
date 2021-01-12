// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"os"
	"strings"
	"time"
)

var (
	LoggerContext = zerolog.
			New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp()
	Logger = LoggerContext.Logger()

	ServiceID   = "hege"
	flagVerbose = 0
	flagQuiet   = false
	flagSyslog  = false
)

func PatchCommandLogs(cmd *cobra.Command) {
	cmd.PersistentFlags().CountVarP(&flagVerbose, "verbose", "v", "Increase the verbosity level")
	cmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", flagQuiet, "Shut the logs")
	cmd.PersistentFlags().BoolVarP(&flagQuiet, "syslog", "s", flagQuiet, "Log in syslog")
	cmd.PersistentFlags().StringVar(&ServiceID, "id", ServiceID, "Use that service ID in the logs")

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		Logger = LoggerContext.Str("id", ServiceID).Logger()

		if flagQuiet {
			zerolog.SetGlobalLevel(zerolog.Disabled)
		} else {
			switch flagVerbose {
			case 0:
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			case 1:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case 2:
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			}
		}
	}
}

func ServerUnaryInterceptorZerolog() grpc.ServerOption {
	return grpc.UnaryInterceptor(newUnaryServerInterceptorZerolog())
}

func ServerStreamInterceptorZerolog() grpc.ServerOption {
	return grpc.StreamInterceptor(newStreamServerInterceptorZerolog())
}

func newStreamServerInterceptorZerolog() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		ctx := ss.Context()
		z := Logger.Info().
			Str("uri", info.FullMethod).
			TimeDiff("t", time.Now(), start)
		if err != nil {
			z.Int("rc", 500)
			z.Err(err)
		} else {
			z.Int("rc", 200)
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			auth := md.Get(":authority")
			if len(auth) > 0 {
				z.Str("local", auth[0])
			}
			sessionID := md.Get("session-id")
			if len(sessionID) > 0 {
				z.Str("session", sessionID[0])
			}
		}
		if peer, ok := peer.FromContext(ctx); ok {
			addr := peer.Addr.String()
			if i := strings.LastIndex(addr, ":"); i > -1 {
				addr = addr[:i]
			}
			z.Str("peer", addr)
		}
		z.Send()
		return err
	}
}

func newUnaryServerInterceptorZerolog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		z := Logger.Info().
			Str("uri", info.FullMethod).
			TimeDiff("t", time.Now(), start)
		if err != nil {
			z.Int("rc", 500)
			z.Err(err)
		} else {
			z.Int("rc", 200)
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			auth := md.Get(":authority")
			if len(auth) > 0 {
				z.Str("local", auth[0])
			}
			sessionID := md.Get("session-id")
			if len(sessionID) > 0 {
				z.Str("session", sessionID[0])
			}
		}
		if peer, ok := peer.FromContext(ctx); ok {
			addr := peer.Addr.String()
			if i := strings.LastIndex(addr, ":"); i > -1 {
				addr = addr[:i]
			}
			z.Str("peer", addr)
		}
		z.Send()
		return resp, err
	}
}

// DumpJSON encodes the argument in JSON and writes the output on os.Stdout
func DumpJSON(x interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	return encoder.Encode(x)
}

// StatusJSON encodes a standard error structure and then forwards it to DumpJSON
func StatusJSON(code int, ID, msg string) error {
	type status struct {
		Msg  string
		Code int
		ID   string
	}
	return DumpJSON(status{Msg: "Create", Code: 200, ID: ID})
}
