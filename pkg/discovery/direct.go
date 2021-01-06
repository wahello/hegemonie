// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package discovery

import (
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/utils"
)

var DefaultDiscovery StatelessDiscovery = TestEnv()

type AllOnHost struct {
	endpoint string
}

func TestEnv() StatelessDiscovery {
	return &AllOnHost{"localhost"}
}

func (d *AllOnHost) makeEndpoint(p uint) (string, error) {
	return fmt.Sprintf("%s:%d", d.endpoint, p), nil
}

func (d *AllOnHost) Kratos() (string, error) {
	return d.makeEndpoint(utils.DefaultPortKratos)
}

func (d *AllOnHost) Keto() (string, error) {
	return d.makeEndpoint(utils.DefaultPortKeto)
}

func (d *AllOnHost) Map() (string, error) {
	return d.makeEndpoint(utils.DefaultPortMap)
}

func (d *AllOnHost) Region() (string, error) {
	return d.makeEndpoint(utils.DefaultPortRegion)
}

func (d *AllOnHost) Event() (string, error) {
	return d.makeEndpoint(utils.DefaultPortEvent)
}
