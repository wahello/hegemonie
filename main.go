// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
	"github.com/jfsmig/hegemonie/region"
	"github.com/jfsmig/hegemonie/front"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "Helpers")
	subcommands.Register(subcommands.FlagsCommand(), "Helpers")
	subcommands.Register(subcommands.CommandsCommand(), "Helpers")
	subcommands.Register(&front.FrontService{}, "Actions")
	subcommands.Register(&region.RegionCommand{}, "Actions")
    flag.Parse()

    ctx := context.Background()
    os.Exit(int(subcommands.Execute(ctx)))
}