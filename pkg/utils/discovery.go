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

type singleHost struct {
	endpoint string
}

type singleEndpoint struct {
	endpoint string
}

// TestEnv forwards to SingleHost on localhost
func TestEnv() StatelessDiscovery { return SingleHost("localhost") }

// TestEnv creates a singleHost implementation.
// singleHost is the simplest implementation of a StatelessDiscovery ever.
// It locates all the services on a given host at their default port value.
func SingleHost(h string) StatelessDiscovery { return &singleHost{h} }

// TestEnv creates a singleEndpoint implementation.
// singleHost is the proxyed implementation of a StatelessDiscovery.
// It locates all the services on a given host, all with the same port.
func SingleEndpoint(e string) StatelessDiscovery { return &singleEndpoint{e} }

func (d *singleHost) makeEndpoint(p uint) (string, error) {
	return fmt.Sprintf("%s:%d", d.endpoint, p), nil
}

// Kratos ... see StatelessDiscovery.Kratos
func (d *singleHost) Kratos() (string, error) { return d.makeEndpoint(DefaultPortKratos) }

// Keto ... see StatelessDiscovery.Keto
func (d *singleHost) Keto() (string, error) { return d.makeEndpoint(DefaultPortKeto) }

// Map ... see StatelessDiscovery.Map
func (d *singleHost) Map() (string, error) { return d.makeEndpoint(DefaultPortMap) }

// Region ... see StatelessDiscovery.Region
func (d *singleHost) Region() (string, error) { return d.makeEndpoint(DefaultPortRegion) }

// Event ... see StatelessDiscovery.Event
func (d *singleHost) Event() (string, error) { return d.makeEndpoint(DefaultPortEvent) }

// Kratos ... see StatelessDiscovery.Kratos
func (d *singleEndpoint) Kratos() (string, error) { return d.endpoint, nil }

// Keto ... see StatelessDiscovery.Keto
func (d *singleEndpoint) Keto() (string, error) { return d.endpoint, nil }

// Map ... see StatelessDiscovery.Map
func (d *singleEndpoint) Map() (string, error) { return d.endpoint, nil }

// Region ... see StatelessDiscovery.Region
func (d *singleEndpoint) Region() (string, error) { return d.endpoint, nil }

// Event ... see StatelessDiscovery.Event
func (d *singleEndpoint) Event() (string, error) { return d.endpoint, nil }
