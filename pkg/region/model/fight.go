// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"sort"
)

func (s SetOfFights) Len() int      { return len(s) }
func (s SetOfFights) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SetOfFights) Less(i, j int) bool {
	return s[i].Cell < s[j].Cell || (s[i].Cell == s[j].Cell && s[i].ID < s[j].ID)
}

func (s SetOfFights) First(cell uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Cell >= cell
	})
	return i
}

func (s *SetOfFights) Add(f *Fight) {
	*s = append(*s, f)
	sort.Sort(*s)
}

func (s SetOfFights) SliceByCell(cell uint64) []*Fight {
	start := s.First(cell)
	for end := start; end < len(s); end++ {
		if s[end].Cell != cell {
			return s[start:end]
		}
	}
	return s[start:]
}
