package rest

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/gokit/pkg/web/mv/json"
	"github.com/ReanSn0w/kincong/cmd/server/cr"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

func New(log lgr.L, revision string, baseURL *url.URL, resolver *cr.CR) *Server {
	return &Server{
		appRevision: revision,
		server:      web.New(log),
		baseURL:     baseURL,
		resolver:    resolver,
	}
}

// @title       Kincong Resolver Rest Server
// @version	    debug
// @description	API для получения подсетей на основе IP адресов, доменов и ASN
// @host        localhost
// @schemes     http
// @BasePath    /
type Server struct {
	appRevision string
	server      *web.Server
	baseURL     *url.URL
	resolver    *cr.CR
}

type ResponseError web.Response[string]

func (s *Server) Start(cancel context.CancelCauseFunc, port int) {
	s.server.Run(cancel, port, s.handler())
}

func (s *Server) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)
}

func (s *Server) handler() http.Handler {
	router := chi.NewRouter()

	web.RestAPI(router, web.NewRestConfig(s.appRevision, s.baseURL))

	router.Route("/api", func(r chi.Router) {
		r.
			With(json.Decoder[ResolveRequest]).
			Post("/resolve", s.resolveHandler)

		r.
			With(json.Decoder[DomainRequest]).
			Post("/domain", s.domainHandler)
	})

	return router
}
