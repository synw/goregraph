package httpServer

import (
	"net/http"
	"fmt"
	"time"
	"os"
	"context"
	"path/filepath"
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/synw/terr"
	"github.com/synw/goregraph/db"
	"github.com/synw/goregraph/lib-r/state"
)


type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func InitHttpServer(serve bool) {
	// routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	// routes
	r.Route("/graphql", func(r chi.Router) {
		r.Get("/*", handleQuery)
	})
	// init
	loc := ":8080"
	httpServer := &http.Server{
		Addr: loc,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	    Handler: r,
	}
	state.HttpServer = httpServer
	// static
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	r.FileServer("/static", http.Dir(filesDir))
	// run
	if state.Verbosity > 0 {
		fmt.Println("Starting http server ...")
	}
	Run()
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

// internal methods

func handleQuery(response http.ResponseWriter, request *http.Request) {
	q := request.URL.Query()["query"][0]
	result, tr := db.RunQuery(q)
	if tr != nil {
		fmt.Println(tr.Formatc())
	}
	json_bytes, err := json.Marshal(result.Data)
	if err != nil {
		fmt.Println(err)
	}
	response = headers(response)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func headers(response http.ResponseWriter) http.ResponseWriter {
	response.Header().Set("Content-Type", "application/json")
	return response
}
