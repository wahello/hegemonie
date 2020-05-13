// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"log"
	"os"
	"text/template"
	"time"
)

var headerTemplate string = `// Code generated : DO NOT EDIT.
// Code generated : {{.Date}}

// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package {{.Package}}

import (
	"errors"
	"sort"
)
`

var bodyTemplate string = `

type {{.SetName}} []{{.TypeName}}

func (s {{.SetName}}) Len() int {
	return len(s)
}

func (s {{.SetName}}) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s {{.SetName}}) Less(i, j int) bool {
	return s[i]{{.Accessor}} < s[j]{{.Accessor}}
}

func (s {{.SetName}}) Check() error {
	if !sort.IsSorted(s) {
		return errors.New("Unsorted")
	}
	var lastId uint64
	for _, a := range s {
		if lastId == a{{.Accessor}} {
			return errors.New("Dupplicate ID")
		}
		lastId = a{{.Accessor}}
	}
	return nil
}

func (s {{.SetName}}) Slice(marker uint64, max uint32) []{{.TypeName}} {
	start := sort.Search(len(s), func(i int) bool {
		return s[i]{{.Accessor}} > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s {{.SetName}}) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i]{{.Accessor}} >= id
	})
	if i < len(s) && s[i]{{.Accessor}} == id {
		return i
	}
	return -1
}

func (s {{.SetName}}) Get(id uint64) {{.TypeName}} {
	var out {{.TypeName}}
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s {{.SetName}}) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *{{.SetName}}) Add(a {{.TypeName}}) {
	*s = append(*s, a)
	if nb := len(*s); nb > 2 && !sort.IsSorted((*s)[nb-2:]) {
		sort.Sort(s)
	}
}

func (s *{{.SetName}}) Remove(a {{.TypeName}}) {
	idx := s.getIndex(a{{.Accessor}})
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

`

type ArrayInstance struct {
	TypeName string
	SetName  string
	Accessor string
	Package  string
	Path     string
	Date     string
}

func main() {
	var err error
	var fout *os.File
	var instance ArrayInstance

	flag.StringVar(&instance.Accessor, "acc", "", "Accessor")
	flag.Parse()

	instance.Package = flag.Arg(0)
	instance.Path = flag.Arg(1)
	instance.TypeName = flag.Arg(2)
	instance.SetName = flag.Arg(3)
	instance.Date = time.Now().String()

	header, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		log.Fatalln("Invalid template", err.Error())
	}
	body, err := template.New("body").Parse(bodyTemplate)
	if err != nil {
		log.Fatalln("Invalid template", err.Error())
	}

	_, err = os.Stat(instance.Path)
	if err == nil {
		fout, err = os.OpenFile(instance.Path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("Invalid output file", err.Error())
		}
	} else {
		fout, err = os.Create(instance.Path)
		if err != nil {
			log.Fatalln("Invalid output file", err.Error())
		}
		header.Execute(fout, &instance)
	}
	defer fout.Close()

	body.Execute(fout, &instance)
}
