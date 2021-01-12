// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import "github.com/juju/errors"

var (
	errRegionExists       = errors.AlreadyExistsf("A region exists with this name")
	errCityExists         = errors.AlreadyExistsf("City exists at that location")
	errCityNotFound       = errors.NotFoundf("No such City")
	errForbidden          = errors.Forbiddenf("Insufficient permissions")
	errNotImplemented     = errors.NotImplementedf("NYI")
	ErrNoSuchUnit         = errors.NotFoundf("No such Unit")
	ErrNotEnoughResources = errors.New("Not enough resources")
)
