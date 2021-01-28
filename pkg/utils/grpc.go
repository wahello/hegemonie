// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"os"
)

// ActionFunc names the signature of the hook to be called when the gRPC connection
// has been established.
type ActionFunc func(ctx context.Context, cli *grpc.ClientConn) error

// Connect establishes a connection to the given service and then call the action.
func Connect(ctx context.Context, endpoint string, action ActionFunc) error {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	cnx, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(grpc_prometheus.UnaryClientInterceptor)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_prometheus.StreamClientInterceptor)),
	)
	if err != nil {
		return errors.Trace(err)
	}
	defer cnx.Close()
	return action(ctx, cnx)
}

// RecvFunc names the signature of the hook that consumes an input and returns an
// object for each PDU received. RecvFunc is initially intended to map any object
// from a gRPC stream into a generic interface{} that is still JSON-encodable.
type RecvFunc func() (interface{}, error)

// EncodeWhole builds an array of all the objects and encodes it in JSON at once.
// Warning: the whole stream will be buffered!
func EncodeWhole(recv RecvFunc) error {
	var out []interface{}
	for {
		x, err := recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Annotate(err, "json encoding error")
		}
		out = append(out, x)
	}
	return DumpJSON(out)
}

// EncodeStream dumps a JSON stream where each line is a JSON-encoded object.
// EncodeStream is initially intended to dump a gRPC stream in JSON at os.Stdout.
func EncodeStream(recv RecvFunc) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "")
	for {
		x, err := recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = encoder.Encode(x)
		if err != nil {
			return errors.Annotate(err, "json encoding error")
		}
	}
	return nil
}
