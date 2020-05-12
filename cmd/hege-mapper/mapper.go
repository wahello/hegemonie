// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
	"math"
	"math/rand"
	"os"
)

type SiteRaw struct {
	Id   string
	X, Y float64
	City bool
}

type RoadRaw struct {
	Src, Dst string
}

type Road struct {
	Src, Dst *Site
}

type MapRaw struct {
	Sites []SiteRaw
	Roads []RoadRaw
}

type Map struct {
	sites map[string]*Site
}

type Site struct {
	raw   SiteRaw
	peers map[*Site]bool
}

func makeMap() Map {
	return Map{
		sites: make(map[string]*Site),
	}
}

func makeRawMap() MapRaw {
	return MapRaw{
		Sites: make([]SiteRaw, 0),
		Roads: make([]RoadRaw, 0),
	}
}

func makeSite(raw SiteRaw) *Site {
	return &Site{
		raw:   raw,
		peers: make(map[*Site]bool),
	}
}

func (mr *MapRaw) Finalize() (Map, error) {
	var err error
	m := makeMap()
	for _, s := range mr.Sites {
		m.sites[s.Id] = &Site{
			raw:   s,
			peers: make(map[*Site]bool),
		}
	}
	for _, r := range mr.Roads {
		if src, ok := m.sites[r.Src]; !ok {
			err = errors.New(fmt.Sprintf("No such site [%s]", r.Src))
			break
		} else if dst, ok := m.sites[r.Dst]; !ok {
			err = errors.New(fmt.Sprintf("No such site [%s]", r.Dst))
			break
		} else {
			src.peers[dst] = true
			dst.peers[src] = true
		}
	}
	return m, err
}

func (s *Site) DotName() string {
	if s.raw.City {
		return s.raw.Id
	} else {
		return "r" + s.raw.Id
	}
}

func (r *Road) Raw() RoadRaw {
	return RoadRaw{Src: r.Src.raw.Id, Dst: r.Dst.raw.Id}
}

func (m *Map) Debug() {
	for _, s := range m.sites {
		log.Println(s.raw)
		for peer, _ := range s.peers {
			log.Println("  ->", peer.raw)
		}
	}
}

func (m *Map) UniqueRoads() <-chan Road {
	out := make(chan Road)
	go func() {
		seen := make(map[RoadRaw]bool)
		for _, s := range m.sites {
			for peer, _ := range s.peers {
				r0 := RoadRaw{Src: s.raw.Id, Dst: peer.raw.Id}
				r1 := RoadRaw{Src: peer.raw.Id, Dst: s.raw.Id}
				if !seen[r0] && !seen[r1] {
					seen[r0] = true
					seen[r1] = true
					out <- Road{s, peer}
				}
			}
		}
		close(out)
	}()
	return out
}

func (m *Map) Raw() MapRaw {
	rm := makeRawMap()
	for _, s := range m.sites {
		rm.Sites = append(rm.Sites, s.raw)
	}
	for r := range m.UniqueRoads() {
		rm.Roads = append(rm.Roads, r.Raw())
	}
	return rm
}

func (m0 *Map) DeepCopy() Map {
	m := makeMap()
	for id, site := range m0.sites {
		m.sites[id] = makeSite(site.raw)
	}
	for _, s := range m0.sites {
		src := m.sites[s.raw.Id]
		for d, _ := range s.peers {
			dst := m.sites[d.raw.Id]
			src.peers[dst] = true
			dst.peers[src] = true
		}
	}
	return m
}

func (m *Map) ComputeBox() (xmin, xmax, ymin, ymax float64) {
	const Max = math.MaxFloat64
	const Min = -Max
	xmin, ymin = Max, Max
	xmax, ymax = Min, Min
	for _, s := range m.sites {
		x, y := s.raw.X, s.raw.Y
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
	}
	if xmin == Max {
		xmin, xmax, ymin, ymax = 0, 0, 0, 0
	}
	return
}

func (m *Map) ShiftAt(xabs, yabs float64) {
	xmin, _, ymin, _ := m.ComputeBox()
	m.Shift(xabs-xmin, yabs-ymin)
}

func (m *Map) Shift(xrel, yrel float64) {
	for _, s := range m.sites {
		s.raw.X += xrel
		s.raw.Y += yrel
	}
}

func (m *Map) ResizeRatio(xratio, yratio float64) {
	for _, s := range m.sites {
		s.raw.X *= xratio
		s.raw.Y *= yratio
	}
}

func (m *Map) ResizeStretch(x, y float64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	m.ResizeRatio(x/xmax, y/ymax)
}

func (m *Map) ResizeAdjust(x, y float64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	xRatio := x / xmax
	yRatio := y / ymax
	ratio := math.Min(xRatio, yRatio)
	m.ResizeRatio(ratio, ratio)
}

func (m *Map) Center(xbound, ybound float64) {
	xmin, xmax, ymin, ymax := m.ComputeBox()
	xdelta, ydelta := xbound-(xmax-xmin), ybound-(ymax-ymin)
	xpad, ypad := xdelta/2.0, ydelta/2.0
	m.Shift(xpad-xmin, ypad-ymin)
}

func (m *Map) splitOneRoad(src, dst *Site, nbSegments uint) {
	if nbSegments < 2 {
		panic("bug")
	}

	xinc := (dst.raw.X - src.raw.X) / float64(nbSegments)
	yinc := (dst.raw.Y - src.raw.Y) / float64(nbSegments)
	segments := make([]*Site, 0, nbSegments+1)

	delete(src.peers, dst)
	delete(dst.peers, src)

	// Create segment boundaries
	segments = append(segments, src)
	for i := uint(0); i < nbSegments-1; i++ {
		id := "x" + uuid.New().String()
		last := segments[len(segments)-1]
		raw := SiteRaw{
			Id:   id,
			City: false,
			X:    last.raw.X + xinc,
			Y:    last.raw.Y + yinc,
		}
		middle := makeSite(raw)
		m.sites[middle.raw.Id] = middle
		segments = append(segments, middle)
	}
	segments = append(segments, dst)

	// Link the segment boundaries
	for i, end := range segments[1:] {
		start := segments[i]
		start.peers[end] = true
		end.peers[start] = true
	}
}

func (m0 *Map) SplitLongRoads(max float64) Map {
	// Work on a deep copy to iterate on the original map while we alter the copy
	m := m0.DeepCopy()
	for r := range m0.UniqueRoads() {
		src := m.sites[r.Src.raw.Id]
		dst := m.sites[r.Dst.raw.Id]
		dist := distance(src, dst)
		if max < dist {
			m.splitOneRoad(src, dst, uint(math.Ceil(dist/max)))
		}
	}
	return m
}

func (m *Map) Noise(xjitter, yjitter float64) {
	for _, s := range m.sites {
		if s.raw.City {
			continue
		}
		s.raw.X += (0.5 - rand.Float64()) * xjitter
		s.raw.Y += (0.5 - rand.Float64()) * yjitter
	}
}

func distance(src, dst *Site) float64 {
	dx := (dst.raw.X - src.raw.X)
	dy := (dst.raw.Y - src.raw.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func CommandNormalize() *cobra.Command {
	return &cobra.Command{
		Use:     "norm",
		Aliases: []string{},
		Short:   "Normalize a map",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)
			encoder := json.NewEncoder(os.Stdout)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			raw = m.Raw()
			return encoder.Encode(&raw)
		},
	}
}

func CommandSplit() *cobra.Command {
	var repeat int = 1
	cmd := &cobra.Command{
		Use:     "split",
		Aliases: []string{},
		Short:   "Split a map",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", " ")

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			for i := 0; i < repeat; i++ {
				m = m.SplitLongRoads(60)
			}
			xmin, xmax, ymin, ymax := m.ComputeBox()
			m.Noise((xmax-xmin)*0.015, (ymax-ymin)*0.015)

			raw = m.Raw()
			return encoder.Encode(&raw)
		},
	}
	cmd.Flags().IntVarP(&repeat, "iterations", "n", 1, "Repeat the split <N> times")
	return cmd
}

func CommandDot() *cobra.Command {
	return &cobra.Command{
		Use:     "dot",
		Aliases: []string{},
		Short:   "Convert the JSON map to DOT (graphviz)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			fmt.Println("graph g {")
			for r := range m.UniqueRoads() {
				fmt.Printf("%s -- %s;\n", r.Src.DotName(), r.Dst.DotName())
			}
			fmt.Println("}")
			return nil
		},
	}
}

func CommandSvg() *cobra.Command {
	var flagStandalone bool

	cmd := &cobra.Command{
		Use:     "svg",
		Aliases: []string{},
		Short:   "Convert the JSON map to SVG",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			xbound, ybound := 1024.0, 768.0
			xPad, yPad := 50.0, 50.0
			m.ResizeAdjust(xbound-2*xPad, ybound-2*yPad)
			m.Center(xbound, ybound)

			fmt.Printf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
	style="background-color: rgb(255, 255, 255);"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	version="1.1"
	width="%dpx" height="%dpx"
	viewBox="-0.5 -0.5 %d %d">
`, int64(xbound), int64(ybound), int64(xbound), int64(ybound))
			fmt.Println(`<g>`)
			for r := range m.UniqueRoads() {
				fmt.Printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="1"/>
`, int64(r.Src.raw.X), int64(r.Src.raw.Y), int64(r.Dst.raw.X), int64(r.Dst.raw.Y))
			}
			fmt.Println(`</g>`)
			fmt.Println(`<g>`)
			for _, s := range m.sites {
				color := `white`
				radius := 5
				stroke := 1
				if s.raw.City {
					color = `gray`
					radius = 10
					stroke = 1
				}
				fmt.Printf(`<circle cx="%d" cy="%d" r="%d" stroke="black" stroke-width="%d" fill="%s"><title>%s</title></circle>
`, int64(s.raw.X), int64(s.raw.Y), radius, stroke, color, s.raw.Id)
				/*
								if s.raw.City {
									fmt.Printf(`<text x="%d" y="%d" fill="#888">%s</text>
				`, int64(s.raw.X) + 10, int64(s.raw.Y) - 10, s.raw.Id)
								}
				*/
			}
			fmt.Println(`</g>`)
			fmt.Println(`</svg>`)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&flagStandalone, "standalone", "1", false, "Also generate the xml header")
	return cmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "mapper",
		Short: "Handle map graphs",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Subcommand required")
		},
	}
	rootCmd.AddCommand(CommandNormalize())
	rootCmd.AddCommand(CommandSplit())
	rootCmd.AddCommand(CommandDot())
	rootCmd.AddCommand(CommandSvg())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
