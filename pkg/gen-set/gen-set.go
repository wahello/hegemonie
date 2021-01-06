// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

var headerTemplate = `// Code generated : DO NOT EDIT.
// Code generated : {{.Date}}

// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package {{.Package}}

import (
	"sort"
	"errors"
)

`

var bodyTemplateCommon = `

type {{.SetName}} []{{.ItemType}}

func (s {{.SetName}}) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s {{.SetName}}) Len() int {
	return len(s)
}

func (s {{.SetName}}) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *{{.SetName}}) Add(a {{.ItemType}}) {
	*s = append(*s, a)
	if nb := len(*s); nb > 2 && !sort.IsSorted((*s)[nb-2:]) {
		sort.Sort(s)
	}
}
`

var bodyTemplate1 = `

func (s {{.SetName}}) Less(i, j int) bool {
	return s[i]{{.F0}} < s[j]{{.F0}}
}

func (s {{.SetName}}) Check() error {
	if !sort.IsSorted(s) {	
		return errors.New("Unsorted")
	}
	var lastId {{.T0}}
	for _, a := range s {
		if lastId == a{{.F0}} {
			return errors.New("Duplicate ID")
		}
		lastId = a{{.F0}}
	}
	return nil
}

func (s {{.SetName}}) Slice(marker {{.T0}}, max uint32) []{{.ItemType}} {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} > marker
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

func (s {{.SetName}}) getIndex(id {{.T0}}) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} >= id
	})
	if i < len(s) && s[i]{{.F0}} == id {
		return i
	}
	return -1
}

func (s {{.SetName}}) Get(id {{.T0}}) {{.ItemType}} {
	var out {{.ItemType}}
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s {{.SetName}}) Has(id {{.T0}}) bool {
	return s.getIndex(id) >= 0
}

func (s *{{.SetName}}) Remove(a {{.ItemType}}) {
	idx := s.getIndex(a{{.F0}})
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

var bodyTemplate2 = `

func (s {{.SetName}}) Less(i, j int) bool {
	p0, p1 := s[i], s[j]
	return p0{{.F0}} < p1{{.F0}} || (p0{{.F0}} == p1{{.F0}} && p0{{.F1}} < p1{{.F1}})
}

func (s {{.SetName}}) First(at {{.T0}}) int {
	return sort.Search(len(s), func(i int) bool { return s[i]{{.F0}} >= at })
}

func (s {{.SetName}}) Check() error {
	if !sort.IsSorted(s) {
		return errors.New("Unsorted")
	}
	var l0 {{.T0}}
	var l1 {{.T1}}
	for _, a := range s {
		if l0 == a{{.F0}} && l1 == a{{.F1}} {
			return errors.New("Duplicate ID")
		}
		l0 = a{{.F0}}
	}
	return nil
}

func (s {{.SetName}}) Slice(m0 {{.T0}}, m1 {{.T1}}, max uint32) []{{.ItemType}} {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}

	iMax := s.Len()
	start := s.First(m0)
	for start < iMax && s[start]{{.F0}} == m0 && s[start]{{.F1}} <= m1 {
		start++
	}

	remaining := uint32(iMax - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s {{.SetName}}) getIndex(f0 {{.T0}}, f1 {{.T1}}) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} >= f0 || (s[i]{{.F0}} == f0 && s[i]{{.F1}} >= f1)
	})
	if i < len(s) && s[i]{{.F0}} == f0 && s[i]{{.F1}} == f1 {
		return i
	}
	return -1
}

func (s {{.SetName}}) Get(f0 {{.T0}}, f1 {{.T1}}) {{.ItemType}} {
	var out {{.ItemType}}
	idx := s.getIndex(f0, f1)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s {{.SetName}}) Has(f0 {{.T0}}, f1 {{.T1}}) bool {
	return s.getIndex(f0, f1) >= 0
}

func (s *{{.SetName}}) Remove(a {{.ItemType}}) {
	idx := s.getIndex(a{{.F0}}, a{{.F1}})
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
	Path string
	Date string

	Package string
	SetName string

	ItemType string

	// First field
	F0 string
	T0 string

	// Second field
	F1 string
	T1 string
}

func main() {
	var err error
	var fout *os.File
	var instance ArrayInstance

	flag.Parse()

	instance.Date = time.Now().String()
	instance.Path = flag.Arg(0)
	instance.F0 = ".ID"
	instance.T0 = "uint64"

	p := flag.Arg(1)
	tokens := strings.Split(p, ":")

	instance.Package = tokens[0]
	instance.SetName = tokens[1]
	instance.ItemType = tokens[2]

	if flag.NArg() > 2 {
		p = flag.Arg(2)
		tokens = strings.Split(p, ":")
		instance.F0 = tokens[0]
		instance.T0 = tokens[1]
		if instance.F0 != "" {
			instance.F0 = "." + instance.F0
		}
		if instance.T0 == "" {
			instance.T0 = "uint64"
		}
	}

	if flag.NArg() > 3 {
		p = flag.Arg(3)
		tokens = strings.Split(p, ":")
		instance.F1 = tokens[0]
		instance.T1 = tokens[1]
		if instance.F1 != "" {
			instance.F1 = "." + instance.F1
		}
		if instance.T1 == "" {
			instance.T1 = "uint64"
		}
	}

	if flag.NArg() > 4 {
		panic("Too many args")
	}

	header, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		log.Fatalln("Invalid template", err.Error())
	}

	common, err := template.New("common").Parse(bodyTemplateCommon)
	if err != nil {
		log.Fatalln("Invalid template", err.Error())
	}

	tpl := bodyTemplate1
	if flag.NArg() == 4 {
		tpl = bodyTemplate2
	}
	body, err := template.New("body").Parse(tpl)
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

	common.Execute(fout, &instance)
	body.Execute(fout, &instance)
}
