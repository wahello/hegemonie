// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import "errors"

var (
	errInvalidState       = errors.New("Structure not initiated")
	errRegionExists       = errors.New("A region exists with this name")
	errCityExists         = errors.New("City exists at that location")
	errCityNotFound       = errors.New("No such City")
	errForbidden          = errors.New("Insufficient permissions")
	errNotImplemented     = errors.New("NYI")
	ErrNoSuchUnit         = errors.New("No such Unit")
	ErrNotEnoughResources = errors.New("Not enough resources")
)
