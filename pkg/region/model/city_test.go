// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"testing"
)

func TestCity_CreateArmyDefence(t *testing.T) {
	ctx := context.Background()
	_, err := NewWorld(ctx)
	if err != nil {
		t.Fatalf("world instanction error: %v", err)
	}
}
