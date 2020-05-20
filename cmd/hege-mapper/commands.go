// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	auth "github.com/jfsmig/hegemonie/pkg/auth/model"
	region "github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/spf13/cobra"
	"os"
)

func CommandNormalize() *cobra.Command {
	return &cobra.Command{
		Use:     "normalize",
		Aliases: []string{"check", "prepare", "sanitize"},
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

			xbound, ybound := 1024.0, 768.0
			xPad, yPad := 50.0, 50.0
			m.ResizeAdjust(xbound-2*xPad, ybound-2*yPad)
			m.Center(xbound, ybound)

			raw = m.Raw()
			return encoder.Encode(&raw)
		},
	}
}

func CommandSplit() *cobra.Command {
	var maxDist float64
	var noise float64

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

			if maxDist > 0 {
				m = m.SplitLongRoads(maxDist)
			}
			xmin, xmax, ymin, ymax := m.ComputeBox()

			if noise > 0 {
				m.Noise((xmax-xmin)*(noise/100), (ymax-ymin)*(noise/100))
			}

			raw = m.Raw()
			return encoder.Encode(&raw)
		},
	}
	cmd.Flags().Float64VarP(&maxDist, "dist", "d", 60, "Max road length")
	cmd.Flags().Float64VarP(&noise, "noise", "n", 15, "Percent of the image dimension used as max noise variation on non-city nodes positions")
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
			for s := range m.SortedSites() {
				color := `white`
				radius := 5
				stroke := 1
				if s.raw.City {
					color = `gray`
					radius = 10
					stroke = 1
				}
				fmt.Printf(`<circle id="%s" class="clickable" cx="%d" cy="%d" r="%d" stroke="black" stroke-width="%d" fill="%s"/>
`, s.raw.Id, int64(s.raw.X), int64(s.raw.Y), radius, stroke, color)
			}
			fmt.Println(`</g>`)
			fmt.Println(`</svg>`)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&flagStandalone, "standalone", "1", false, "Also generate the xml header")
	return cmd
}

func CommandExport() *cobra.Command {
	var config string

	cmd := &cobra.Command{
		Use:     "export",
		Aliases: []string{"finish"},
		Short:   "Export the map as JSON files as expected by a Region agent",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var dirOut string

			switch len(args) {
			case 0:
				return errors.New("Expected argument: path to the output directory")
			case 1:
				dirOut = args[0]
			default:
				return errors.New("")
			}

			var aaa auth.Db
			aaa.Init()
			aaa.ReHash()
			u := aaa.Create("admin@hegemonie.be")
			u.Name = "Super Admin"
			u.Admin = true
			u.Password = ":plop"
			u.Characters = append(u.Characters, auth.Character{
				Id: aaa.NextId,
				Name: "Waku",
				Region: "Calaquyr",
			})
			aaa.ReHash()
			if err = aaa.Check(); err != nil {
				return err
			}

			var raw MapRaw
			decoder := json.NewDecoder(os.Stdin)
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			w := region.World{}
			w.Init()

			// Load the configuration, because we need models
			w.LoadDefinitionsFromFiles(config)

			// Fill the world with cities and map cells
			site2cell := make(map[*Site]*region.MapVertex)
			for site := range m.SortedSites() {
				cell := w.Places.CellCreate()
				cell.X = uint64(site.raw.X)
				cell.Y = uint64(site.raw.Y)
				if site.raw.City {
					cityId, err := w.CityCreateRandom(cell.Id)
					if err != nil {
						return err
					}
					city := w.CityGet(cityId)
					city.Name = site.raw.Id
					city.Cell = cell.Id
					cell.City = cityId
				}
				site2cell[site] = cell
			}
			for road := range m.UniqueRoads() {
				src := site2cell[road.Src]
				dst := site2cell[road.Dst]
				if err = w.Places.RoadCreate(src.Id, dst.Id, true); err != nil {
					return err
				}
				if err = w.Places.RoadCreate(dst.Id, src.Id, true); err != nil {
					return err
				}
			}

			if err = w.PostLoad(); err != nil {
				return err
			}
			if err = w.Check(); err != nil {
				return err
			}

			// Patch the resource multipliers
			one := region.MultiplierUniform(1.0)
			for _, x := range w.Definitions.Buildings {
				x.Prod.Mult.Add(one)
			}
			for _, x := range w.Definitions.Knowledges {
				x.Prod.Mult.Add(one)
			}
			for _, x := range w.Definitions.Units {
				x.Prod.Mult.Add(one)
			}

			// Populate the cities with a set of minimal troops / units
			for _, pCity := range w.Live.Cities {
				pCity.Owner = u.Characters[0].Id
				// Create one Army per City
				_ = pCity.UnitCreate(&w, w.Definitions.Units[0])
				_ = pCity.MakeDefence(&w)
				// Create one Unit per City
				_ = pCity.UnitCreate(&w, w.Definitions.Units[0])
			}

			// Dump the LIVE base of the world concerned by the current script
			if p, err := w.SaveLiveToFiles(dirOut); err != nil {
				return err
			} else {
				fmt.Println("HEGE_LIVE=" + p)
			}

			// Dump the configuration
			if p, err := w.SaveDefinitionsToFiles(dirOut + "/definitions"); err != nil {
				return err
			} else {
				fmt.Println("HEGE_DEFS=" + p)
			}

			if f, err := os.Create(dirOut + "/auth.json"); err != nil {
				return err
			} else {
				encoder := json.NewEncoder(f)
				encoder.SetIndent("", " ")
				err = encoder.Encode(aaa.UsersById)
				_ = f.Close()
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&config, "config", "c", "", "Configuration Directory used to load the City patterns")
	return cmd
}
