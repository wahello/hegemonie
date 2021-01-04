// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import "fmt"

const (
	// Authentication by ORY
	DefaultPortKratos = 4434
	// Autorisation by ORY
	DefaultPortKeto = 4466
	// OpenID connect by ORY
	DefaultPortHydra = 6686

	// Hegemonie Region internal API
	DefaultPortRegion = 8081
	// Hegemonie Event internal API
	DefaultPortEvent = 8082
	// Hegemonie Map internal API
	DefaultPortMap = 8083
)

func EndpointLocal(port uint) string { return fmt.Sprintf("localhost:%v", port) }
func EndpointAny(port uint) string   { return fmt.Sprintf("0.0.0.0:%v", port) }
