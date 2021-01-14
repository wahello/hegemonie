// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

// Map actions that are exposed to a World
type MapView interface {
	// Step resolves the next step of the from src to dst
	// the context is typically inheritated from the original request context.
	Step(ctx context.Context, mapName string, src, dst uint64) (uint64, error)
}

type directPathResolver struct{}

func newDirectPathResolver(_ context.Context) (MapView, error) {
	return &directPathResolver{}, nil
}

func (r *directPathResolver) Step(ctx context.Context, mapName string, src, dst uint64) (uint64, error) {
	// TODO(jfs): keep a cache of the map connection
	endpoint, err := utils.DefaultDiscovery.Map()
	if err != nil {
		return 0, errors.Annotate(err, "map service not located")
	}
	cnx, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return 0, err
	}
	defer cnx.Close()

	return 0, errors.NotImplementedf("NYI")
}
