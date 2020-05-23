// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/memcache"
	"gopkg.in/macaron.v1"
)

func Command() *cobra.Command {
	front := FrontService{}
	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Web service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if front.endpointRegion == "" {
				return errors.New("Missing region URL")
			}

			if fi, err := os.Stat(front.dirTemplates); err != nil || !fi.IsDir() {
				return errors.New("Invalid path for the directory of templates")
			}
			if fi, err := os.Stat(front.dirStatic); err != nil || !fi.IsDir() {
				return errors.New("Invalid path for the directory of static files")
			}

			m := macaron.New()
			m.Use(macaron.Recovery())
			m.SetDefaultCookieSecret("hegemonie-session-NOT-SET")
			m.Use(pongo2.Pongoer(pongo2.Options{
				Directory:       front.dirTemplates,
				Extensions:      []string{".tpl", ".html", ".tmpl"},
				HTMLContentType: "text/html",
				Charset:         "UTF-8",
				IndentJSON:      true,
				IndentXML:       true,
			}))
			m.Use(session.Sessioner(session.Options{
				Provider:       "memcache",
				ProviderConfig: "127.0.0.1:11211",
			}))
			m.Use(zeroLogger())
			m.Use(func(ctx *macaron.Context, s session.Store) {
				auth := func() {
					uid := s.Get("userid")
					if uid == "" {
						ctx.Redirect("/index.html")
					}
				}
				// Pages under the /game/* prefix require an established authentication
				switch {
				case strings.HasPrefix(ctx.Req.URL.Path, "/game/"),
					strings.HasPrefix(ctx.Req.URL.Path, "/action/"):
					auth()
				}
			})
			front.routePages(m)
			m.Use(macaron.Static(front.dirStatic, macaron.StaticOptions{
				Prefix:      "static",
				SkipLogging: true,
			}))
			front.routeForms(m)

			var err error

			front.cnxAuth, err = grpc.Dial(front.endpointAuth, grpc.WithInsecure())
			if err != nil {
				return err
			}

			front.cnxRegion, err = grpc.Dial(front.endpointRegion, grpc.WithInsecure())
			if err != nil {
				return err
			}

			go front.loopReload(context.Background())

			return http.ListenAndServe(front.endpointNorth, m)
		},
	}
	agent.Flags().StringVar(&front.endpointNorth, "endpoint", ":8080", "TCP/IP North endpoint")
	agent.Flags().StringVar(&front.endpointRegion, "region", "", "World Server to be contacted")
	agent.Flags().StringVar(&front.endpointAuth, "auth", "", "Auth Server to be contacted")
	agent.Flags().StringVar(&front.dirTemplates, "templates", "/data/templates", "Directory with the HTML templates")
	agent.Flags().StringVar(&front.dirStatic, "static", "/data/static", "Directory with the static files")
	return agent
}

type FrontService struct {
	dirTemplates   string
	dirStatic      string
	endpointNorth  string
	endpointRegion string
	endpointAuth   string

	cnxRegion *grpc.ClientConn
	cnxAuth   *grpc.ClientConn

	rw        sync.RWMutex
	units     map[uint64]*region.UnitTypeView
	buildings map[uint64]*region.BuildingTypeView
	knowledge map[uint64]*region.KnowledgeTypeView
}

func (f *FrontService) reload(sessionId string, ctx0 context.Context) {
	cli := region.NewDefinitionsClient(f.cnxRegion)
	ctx := metadata.AppendToOutgoingContext(ctx0, "session-id", sessionId)

	func() {
		last := uint64(0)
		tab := make(map[uint64]*region.UnitTypeView)

		for {
			args := &region.PaginatedQuery{Marker: last, Max: 100}
			l, err := cli.ListUnits(ctx, args)
			if err != nil {
				log.Println("Reload error (units):", err.Error())
				return
			}
			if len(l.Items) <= 0 {
				break
			}
			for _, item := range l.Items {
				if last < item.Id {
					last = item.Id
				}
				tab[item.Id] = item
			}
		}
		if len(tab) > 0 {
			f.rw.Lock()
			f.units = tab
			f.rw.Unlock()
		}
	}()

	func() {
		last := uint64(0)
		tab := make(map[uint64]*region.BuildingTypeView)
		for {
			args := &region.PaginatedQuery{Marker: last, Max: 100}
			l, err := cli.ListBuildings(ctx, args)
			if err != nil {
				log.Println("Reload error (buildings):", err.Error())
				return
			}
			if len(l.Items) <= 0 {
				break
			}
			for _, item := range l.Items {
				if last < item.Id {
					last = item.Id
				}
				tab[item.Id] = item
			}
		}
		if len(tab) > 0 {
			f.rw.Lock()
			f.buildings = tab
			f.rw.Unlock()
		}
	}()

	func() {
		last := uint64(0)
		tab := make(map[uint64]*region.KnowledgeTypeView)
		for {
			args := &region.PaginatedQuery{Marker: last, Max: 100}
			l, err := cli.ListKnowledges(ctx, args)
			if err != nil {
				log.Println("Reload error (knowledge):", err.Error())
				return
			}
			if len(l.Items) <= 0 {
				break
			}
			for _, item := range l.Items {
				if last < item.Id {
					last = item.Id
				}
				tab[item.Id] = item
			}
		}
		if len(tab) > 0 {
			f.rw.Lock()
			f.knowledge = tab
			f.rw.Unlock()
		}
	}()
}

func (f *FrontService) loopReload(ctx context.Context) {
	go func() {
		sessionId := uuid.New().String()
		for _, v := range []int{2, 4, 8, 16} {
			f.reload(sessionId, ctx)
			<-time.After(time.Duration(v) * time.Second)
		}
		for {
			f.reload(sessionId, ctx)
			<-time.After(61 * time.Second)
		}
	}()
}

func utoa(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func atou(s string) uint64 {
	u, err := strconv.ParseUint(s, 10, 63)
	if err != nil {
		return 0
	} else {
		return u
	}
}

func ptou(p interface{}) uint64 {
	if p == nil {
		return 0
	}
	return atou(p.(string))
}

func zeroLogger() macaron.Handler {
	return func(ctx *macaron.Context, s session.Store) {
		start := time.Now()
		rw := ctx.Resp.(macaron.ResponseWriter)
		ctx.Next()
		z := utils.Logger.Info().
			Str("peer", ctx.RemoteAddr()).
			Str("local", ctx.Req.Host).
			Str("verb", ctx.Req.Method).
			Str("uri", ctx.Req.RequestURI).
			Int("rc", rw.Status()).
			TimeDiff("t", time.Now(), start)
		if sessionId := s.Get("session-id"); sessionId != nil {
			z.Str("session", sessionId.(string))
		}
		z.Send()
	}
}

func contextMacaronToGrpc(ctx *macaron.Context, s session.Store) context.Context {
	return metadata.AppendToOutgoingContext(ctx.Req.Context(),
		"session-id", s.Get("session-id").(string))
}
