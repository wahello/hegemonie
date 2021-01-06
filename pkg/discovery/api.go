// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package discovery

// Interface to be implemented to act as a service source in Hegemonie.
type StatelessDiscovery interface {
	// Inform about the ORY Kratos services (Authentication)
	Kratos() (string, error)

	// Inform about the ORY Keto services (Authorisation)
	Keto() (string, error)

	// Inform about hegemonie's Map services
	Map() (string, error)

	// Inform about hegemonie's Region services
	// Please note that those services are typically sharded. Stateless weighted polling
	// is only meaningful when it is necessary to instantiate a new Region.
	Region() (string, error)

	// Inform about hegemonie's Event services
	Event() (string, error)
}
