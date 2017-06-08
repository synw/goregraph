package server

import (
	"context"
	"fmt"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	g "github.com/synw/goregraph/lib-r/httpServer"
	"github.com/synw/microb-goregraph/datatypes"
	"github.com/synw/microb-goregraph/state"
	"github.com/synw/terr"
	"net/http"
	"time"
)

var Server = &datatypes.GraphqlServer{
	false,
	&http.Server{},
}

func InitServer() {
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	// routes
	r.Route("/graphql", func(r chi.Router) {
		r.Get("/*", g.HandleQuery)
	})
	// init
	host := state.Conf.Host
	httpServer := &http.Server{
		Addr:         host,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	// set server
	Server = &datatypes.GraphqlServer{
		false,
		httpServer,
	}
	Run()
}

func Stop() *terr.Trace {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := Server.Instance
	err := srv.Shutdown(ctx)
	if err != nil {
		tr := terr.New("server.Stop", err)
		return tr
	}
	Server.Running = false
	return nil
}

func Run() {
	Server.Running = true
	if state.Verbosity > 0 {
		host := state.Conf.Host
		fmt.Println("Graphql server is up at " + host + "...")
	}
	Server.Instance.ListenAndServe()
}
