package httpServer

import (
	"net/http"
	"fmt"
	"time"
	"context"
	"errors"
	"encoding/json"
	"github.com/goware/cors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/graphql-go/graphql"
	"github.com/synw/terr"
	"github.com/synw/goregraph/db"
	"github.com/synw/goregraph/lib-r/state"
)


type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func InitHttpServer(serve bool) {
	r := chi.NewRouter()
	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	cors := cors.New(cors.Options{
		AllowedOrigins: state.Conf.Cors,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	})
	r.Use(cors.Handler)
  
	// routes
	r.Route("/graphql", func(r chi.Router) {
		r.Get("/*", HandleQuery)
	})
	// init
	httpServer := &http.Server{
		Addr: state.Addr,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: r,
	}
	state.HttpServer = httpServer
	// run
	if state.Verbosity > 0 {
		fmt.Println("Starting http server ...")
	}
	if serve { 
		Run() 
	}
}

func Run() {
	state.HttpServer.ListenAndServe()
}

func Stop() *terr.Trace {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := state.HttpServer
	err := srv.Shutdown(ctx)
	if err != nil {
		tr := terr.New("httpServer.Stop", err)
		return tr
	}
	return nil
}

func HandleQuery(response http.ResponseWriter, request *http.Request) {
	q := request.URL.Query()["query"][0]
	res := graphql.Do(graphql.Params{
		Schema: db.Schem,
		RequestString: q,
	})
	if len(res.Errors) > 0 {
		msg := fmt.Sprintf("wrong res, unexpected errors: %v", res.Errors)
		err := errors.New(msg)
		tr := terr.New("httpServer.handleQuery", err)
		tr.Printf("httpServer.handleQuery")
	}
	data := res.Data
	json_bytes, _ := json.Marshal(data)
	response = headers(response)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

// internal methods

func headers(response http.ResponseWriter) http.ResponseWriter {
	response.Header().Set("Content-Type", "application/json")
	return response
}
