// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
	"time"
)

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
