// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
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
	creds := credentials.NewTLS(config)

	cnx, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(creds))
	if err != nil {
		return errors.Trace(err)
	}
	defer cnx.Close()
	return action(ctx, cnx)
}

// ServeTLS automates the creation of a grpc.Server over a TLS connection
// with the proper interceptors.
func ServerTLS(pathKey, pathCrt string) (*grpc.Server, error) {
	if len(pathCrt) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 certificate path [%s]", pathCrt)
	}
	if len(pathKey) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 key path [%s]", pathKey)
	}
	var certBytes, keyBytes []byte
	var err error

	Logger.Info().Str("key", pathKey).Str("crt", pathCrt).Msg("TLS config")

	if certBytes, err = ioutil.ReadFile(pathCrt); err != nil {
		return nil, errors.Annotate(err, "certificate file error")
	}
	if keyBytes, err = ioutil.ReadFile(pathKey); err != nil {
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

	creds := credentials.NewServerTLSFromCert(&cert)
	srv := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			newUnaryServerInterceptorZerolog())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			newStreamServerInterceptorZerolog())))
	return srv, nil
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
