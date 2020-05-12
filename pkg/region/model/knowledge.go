// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

// TODO(jfs): Maybe speed the execution with a reverse index of Requires
func (s SetOfKnowledgeTypes) Frontier(owned []*Knowledge) []*KnowledgeType {
	pending := make(map[uint64]bool)
	finished := make(map[uint64]bool)
	for _, k := range owned {
		if k.Ticks == 0 {
			finished[k.Type] = true
		} else {
			pending[k.Type] = true
		}
	}

	valid := func(kt *KnowledgeType) bool {
		if finished[kt.Id] || pending[kt.Id] {
			return false
		}
		for _, c := range kt.Conflicts {
			if finished[c] || pending[c] {
				return false
			}
		}
		for _, c := range kt.Requires {
			if !finished[c] {
				return false
			}
		}
		return true
	}

	result := make([]*KnowledgeType, 0)
	for _, kt := range s {
		if valid(kt) {
			result = append(result, kt)
		}
	}
	return result
}
