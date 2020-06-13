// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"errors"
	"github.com/go-macaron/binding"
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
			m.SetDefaultCookieSecret("heged-session-NOT-SET")
			m.Use(macaron.Recovery())
			m.Use(macaron.Static(front.dirStatic, macaron.StaticOptions{
				Prefix:      "static",
				SkipLogging: true,
			}))
			m.Use(session.Sessioner(session.Options{
				Provider:       "memcache",
				ProviderConfig: "127.0.0.1:11211",
			}))
			m.Use(zeroLogger())
			m.Use(pongo2.Pongoer(pongo2.Options{
				Directory:       front.dirTemplates,
				Extensions:      []string{".tpl", ".html", ".tmpl"},
				HTMLContentType: "text/html",
				Charset:         "UTF-8",
				IndentJSON:      true,
				IndentXML:       true,
			}))
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
			m.Post("/action/login", binding.Bind(FormLogin{}), doLogin(&front, m))
			m.Post("/action/logout", doLogout(&front, m))
			m.Get("/action/logout", doLogout(&front, m))
			m.Post("/action/move", doMove(&front, m))
			m.Post("/action/produce", doProduce(&front, m))
			m.Post("/action/city/study", binding.Bind(FormCityStudy{}), doCityStudy(&front, m))
			m.Post("/action/city/build", binding.Bind(FormCityBuild{}), doCityBuild(&front, m))
			m.Post("/action/city/train", binding.Bind(FormCityTrain{}), doCityTrain(&front, m))
			m.Post("/action/army/cancel", binding.Bind(FormCityArmyDisband{}), doCityCancelArmy(&front, m))
			m.Post("/action/army/disband", binding.Bind(FormCityArmyDisband{}), doCityDisbandArmy(&front, m))
			m.Post("/action/army/command", binding.Bind(FormCityArmyCommand{}), doCityCommandArmy(&front, m))
			m.Post("/action/army/create", binding.Bind(FormCityArmyCreate{}), doCityCreateArmy(&front, m))
			m.Post("/action/city/unit/transfer", binding.Bind(FormCityUnitTransfer{}), doCityTransferUnit(&front, m))
			m.Get("/game/admin", serveGameAdmin(&front))
			m.Get("/game/user", serveGameUser(&front))
			m.Get("/game/character", serveGameCharacter(&front))
			m.Get("/game/land/overview", serveGameCityOverview(&front))
			m.Get("/game/land/budget", serveGameCityBudget(&front))
			m.Get("/game/land/buildings", serveGameCityBuildings(&front))
			m.Get("/game/land/armies", serveGameCityArmies(&front))
			m.Get("/game/land/units", serveGameCityUnits(&front))
			m.Get("/game/land/knowledges", serveGameCityKnowledges(&front))
			m.Get("/game/army", serveGameArmyDetail(&front))
			m.Get("/map/region", serveRegionMap(&front))
			m.Get("/map/cities", serveRegionCities(&front))
			m.Get("/", serveRoot)

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
	cities    map[uint64]*region.PublicCity
	locations map[uint64]*region.Vertex
}

func (f *FrontService) loadAllCities(ctx context.Context, cli region.MapClient) (map[uint64]*region.PublicCity, error) {
	last := uint64(0)
	tab := make(map[uint64]*region.PublicCity)

	for {
		args := &region.PaginatedQuery{Marker: last, Max: 1000}
		l, err := cli.Cities(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if last < item.Id {
				last = item.Id
			}
			tab[item.Id] = item
		}
	}
}

func (f *FrontService) loadAllLocations(ctx context.Context, cli region.MapClient) (map[uint64]*region.Vertex, error) {
	last := uint64(0)
	tab := make(map[uint64]*region.Vertex)

	for {
		args := &region.PaginatedQuery{Marker: last, Max: 10000}
		l, err := cli.Vertices(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if last < item.Id {
				last = item.Id
			}
			tab[item.Id] = item
		}
	}
}

func (f *FrontService) loadAllRoads(ctx context.Context, cli region.MapClient) ([]*region.Edge, error) {
	var lastSrc, lastDst uint64
	tab := make([]*region.Edge, 0)

	for {
		args := &region.ListEdgesReq{MarkerSrc: lastSrc, MarkerDst: lastDst, Max: 10000}
		l, err := cli.Edges(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if lastSrc < item.Src {
				lastSrc = item.Src
				lastDst = item.Dst
			} else if lastSrc == item.Src && lastDst < item.Dst {
				lastDst = item.Dst
			}
			tab = append(tab, item)
		}
	}
}

func (f *FrontService) loadAllUnits(ctx context.Context, cli region.DefinitionsClient) (map[uint64]*region.UnitTypeView, error) {
	last := uint64(0)
	tab := make(map[uint64]*region.UnitTypeView)

	for {
		args := &region.PaginatedQuery{Marker: last, Max: 1000}
		l, err := cli.ListUnits(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if last < item.Id {
				last = item.Id
			}
			tab[item.Id] = item
		}
	}
}

func (f *FrontService) loadAllBuildings(ctx context.Context, cli region.DefinitionsClient) (map[uint64]*region.BuildingTypeView, error) {
	last := uint64(0)
	tab := make(map[uint64]*region.BuildingTypeView)

	for {
		args := &region.PaginatedQuery{Marker: last, Max: 1000}
		l, err := cli.ListBuildings(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if last < item.Id {
				last = item.Id
			}
			tab[item.Id] = item
		}
	}
}

func (f *FrontService) loadAllKnowledges(ctx context.Context, cli region.DefinitionsClient) (map[uint64]*region.KnowledgeTypeView, error) {
	last := uint64(0)
	tab := make(map[uint64]*region.KnowledgeTypeView)

	for {
		args := &region.PaginatedQuery{Marker: last, Max: 1000}
		l, err := cli.ListKnowledges(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(l.Items) <= 0 {
			return tab, nil
		}
		for _, item := range l.Items {
			if last < item.Id {
				last = item.Id
			}
			tab[item.Id] = item
		}
	}
}

func (f *FrontService) reload(cli region.DefinitionsClient, sessionId string, ctx0 context.Context) {
	ctx := metadata.AppendToOutgoingContext(ctx0, "session-id", sessionId)

	var uerr, berr, kerr error
	var wg sync.WaitGroup
	var utv map[uint64]*region.UnitTypeView
	var btv map[uint64]*region.BuildingTypeView
	var ktv map[uint64]*region.KnowledgeTypeView

	wg.Add(3)
	go func() {
		defer wg.Done()
		utv, uerr = f.loadAllUnits(ctx, cli)
	}()
	go func() {
		defer wg.Done()
		btv, berr = f.loadAllBuildings(ctx, cli)
	}()
	go func() {
		defer wg.Done()
		ktv, kerr = f.loadAllKnowledges(ctx, cli)
	}()
	wg.Wait()

	if uerr != nil {
		utils.Logger.Warn().Err(uerr).Str("step", "units").Msg("Reload error")
	}
	if berr != nil {
		utils.Logger.Warn().Err(kerr).Str("step", "buildings").Msg("Reload error")
	}
	if kerr != nil {
		utils.Logger.Warn().Err(berr).Str("step", "knowledge").Msg("Reload error")
	}

	f.rw.Lock()
	if uerr == nil {
		f.units = utv
	}
	if berr == nil {
		f.buildings = btv
	}
	if kerr == nil {
		f.knowledge = ktv
	}
	f.rw.Unlock()
}

func (f *FrontService) loopReload(ctx context.Context) {
	sessionId := uuid.New().String()
	for _, v := range []int{2, 4, 8, 16} {
		cli := region.NewDefinitionsClient(f.cnxRegion)
		f.reload(cli, sessionId, ctx)
		<-time.After(time.Duration(v) * time.Second)
	}
	for {
		cli := region.NewDefinitionsClient(f.cnxRegion)
		f.reload(cli, sessionId, ctx)
		<-time.After(61 * time.Second)
	}
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
	return contextPatchToGrpc(ctx.Req.Context(), s)
}

func contextPatchToGrpc(ctx context.Context, s session.Store) context.Context {
	return contextSessionToGrpc(ctx, s.Get("session-id").(string))
}

func contextSessionToGrpc(ctx context.Context, sessionId string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "session-id", sessionId)
}
