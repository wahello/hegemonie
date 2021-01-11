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
	"errors"
	"fmt"
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
		return err
	}
	defer cnx.Close()
	return action(ctx, cnx)
}

// ServeTLS automates the creation of a grpc.Server over a TLS connection
// with the proper interceptors.
func ServerTLS(pathKey, pathCrt string) (*grpc.Server, error) {
	if len(pathCrt) <= 0 {
		return nil, fmt.Errorf("invalid TLS/x509 certificate path [%s]", pathCrt)
	}
	if len(pathKey) <= 0 {
		return nil, fmt.Errorf("invalid TLS/x509 key path [%s]", pathKey)
	}
	var certBytes, keyBytes []byte
	var err error

	Logger.Info().Str("key", pathKey).Str("crt", pathCrt).Msg("TLS config")

	if certBytes, err = ioutil.ReadFile(pathCrt); err != nil {
		return nil, err
	}
	if keyBytes, err = ioutil.ReadFile(pathKey); err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, errors.New("invalid certificates")
	}

	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewServerTLSFromCert(&cert)
	srv := grpc.NewServer(
		grpc.Creds(creds),
		grpc.StreamInterceptor(newStreamServerInterceptorZerolog()),
		grpc.UnaryInterceptor(newUnaryServerInterceptorZerolog()))
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
			return err
		}
		out = append(out, x)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
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
			return err
		}
	}
	return nil
}