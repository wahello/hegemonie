// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"fmt"
)

// StatelessDiscovery is the simplest form of Discovery API providing one call per
// type of service. Each call returns either a usable endpoint string or the error
// that occurred during the discovery process.
// The implementation of the StatelessDiscovery interface is responsible for the
// management of its concurrent accesses.
type StatelessDiscovery interface {
	// Kratos locates an ORY kratos service (Authentication)
	Kratos() (string, error)

	// Keto locates an ORY keto service (Authorisation)
	Keto() (string, error)

	// Map locates map services in Hegemonie
	Map() (string, error)

	// Region locates an hegemonie's region service
	// Please note that those services are typically sharded. Stateless weighted polling
	// is only meaningful when it is necessary to instantiate a new Region.
	Region() (string, error)

	// Event locates an hegemonie's event services
	Event() (string, error)
}

// DefaultDiscovery is the default implementation of a discovery.
// Valued by default to the discovery of test services, all located on
// localhost and serving default ports.
var DefaultDiscovery StatelessDiscovery = TestEnv()

// AllOnHost is the simplest implementation of a StatelessDiscovery ever.
// It locates all the services on a given host at their default port value.
type AllOnHost struct {
	endpoint string
}

// TestEnv creates a AllOnHost implementation based on localhost.
func TestEnv() StatelessDiscovery {
	return &AllOnHost{"localhost"}
}

func (d *AllOnHost) makeEndpoint(p uint) (string, error) {
	return fmt.Sprintf("%s:%d", d.endpoint, p), nil
}

// see StatelessDiscovery.Kratos
func (d *AllOnHost) Kratos() (string, error) {
	return d.makeEndpoint(DefaultPortKratos)
}

// see StatelessDiscovery.Keto
func (d *AllOnHost) Keto() (string, error) {
	return d.makeEndpoint(DefaultPortKeto)
}

// see StatelessDiscovery.Map
func (d *AllOnHost) Map() (string, error) {
	return d.makeEndpoint(DefaultPortMap)
}

// see StatelessDiscovery.Region
func (d *AllOnHost) Region() (string, error) {
	return d.makeEndpoint(DefaultPortRegion)
}

// see StatelessDiscovery.Event
func (d *AllOnHost) Event() (string, error) {
	return d.makeEndpoint(DefaultPortEvent)
}
