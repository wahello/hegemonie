// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

func (s *SetOfKnowledgeTypes) Len() int {
	return len(*s)
}

func (s *SetOfKnowledgeTypes) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfKnowledgeTypes) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfKnowledgeTypes) Add(b *KnowledgeType) {
	*s = append(*s, b)
	sort.Sort(s)
}

func (s *SetOfKnowledgeTypes) Get(id uint64) *KnowledgeType {
	for _, i := range *s {
		if i.Id == id {
			return i
		}
	}
	return nil
}

func (s *SetOfKnowledgeTypes) Has(id uint64) bool {
	return s.Get(id) != nil
}

func (s *SetOfKnowledges) Len() int {
	return len(*s)
}

func (s *SetOfKnowledges) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfKnowledges) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (s *SetOfKnowledges) Add(b *Knowledge) {
	*s = append(*s, b)
	sort.Sort(s)
}

func (w *World) KnowledgeTypeGet(id uint64) *KnowledgeType {
	return w.Definitions.Knowledges.Get(id)
}

// TODO(jfs): Maybe speed the execution with a reverse index of Requires
func (w *World) KnowledgeGetFrontier(owned []*Knowledge) []*KnowledgeType {
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
	for _, kt := range w.Definitions.Knowledges {
		if valid(kt) {
			result = append(result, kt)
		}
	}
	return result
}
