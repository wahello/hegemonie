// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"context"
	"google.golang.org/grpc"
)

type ActionFunc func(ctx context.Context, cli *grpc.ClientConn) error

func Connect(ctx context.Context, endpoint string, action ActionFunc) error {
	cnx, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	return action(ctx, cnx)
}
