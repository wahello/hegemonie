// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type PersistencyMapping []CfgSection

type CfgSection struct {
	Path string
	Obj  interface{}
}

func (p PersistencyMapping) Dump() error {
	for _, section := range p {
		out, err := os.Create(section.Path)
		if err != nil {
			return fmt.Errorf("Failed to save the World in [%s]: %s", section.Path, err.Error())
		}
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", " ")
		err = encoder.Encode(section.Obj)
		_ = out.Close()
		if err != nil {
			return fmt.Errorf("Failed to save the World in [%s]: %s", section.Path, err.Error())
		}
	}
	return nil
}

func (p PersistencyMapping) Load() error {
	for _, section := range p {
		in, err := os.Open(section.Path)
		if err != nil {
			return fmt.Errorf("Failed to load the World from [%s]: %s", section.Path, err.Error())
		}
		err = json.NewDecoder(in).Decode(section.Obj)
		in.Close()
		if err != nil {
			return fmt.Errorf("Failed to load the World from [%s]: %s", section.Path, err.Error())
		}
	}
	return nil
}
