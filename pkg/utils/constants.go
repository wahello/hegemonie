// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

const (
	// DefaultPortKratos is the default port of the Kratos authentication service by ORY
	DefaultPortKratos = 4434
	// DefaultPortKeto is the default port of the Keto autorisation service by ORY
	DefaultPortKeto = 4466
	// DefaultPortHydra is the default port of the Hydra OpenID Connect / OAuth2 provider by ORY
	DefaultPortHydra = 6686

	// DefaultPortCommon is the default port used to hast the service.
	// This is also the configuration expected by the docker image.
	DefaultPortCommon = 6000

	// DefaultPortMonitoring is the default port used for clear-text HTTP
	// serving the /metrics route dedicated to Prometheus exporters.
	DefaultPortMonitoring = 6001

	// DefaultPortRegion is the default port of the Hegemonie Region internal API service
	DefaultPortRegion = 8081
	// DefaultPortEvent is the default port of the Hegemonie Event internal API service
	DefaultPortEvent = 8082
	// DefaultPortMap is the default port of the Hegemonie Map internal API service
	DefaultPortMap = 8083
)
