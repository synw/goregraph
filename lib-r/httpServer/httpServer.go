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
	"reflect"
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
		r.Get("/*", handleQuery)
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

// internal methods

func handleQuery(response http.ResponseWriter, request *http.Request) {
	q := request.URL.Query()["query"][0]
	/*res, tr := db.RunQuery(q)
	if tr != nil {
		fmt.Println(tr.Formatc())
	}*/
	res := graphql.Do(graphql.Params{
		Schema: db.Schema,
		RequestString: q,
	})
	if len(res.Errors) > 0 {
		msg := fmt.Sprintf("wrong res, unexpected errors: %v", res.Errors)
		err := errors.New(msg)
		tr := terr.New("httpServer.handleQuery", err)
		tr.Printf("httpServer.handleQuery")
	}
	//w := json.NewEncoder(response).Encode(res)
	//fmt.Println("RES", res)
	//return w
	
	data := res.Data
	
	//json_bytes := data.Obj.Bytes()
	fmt.Println("DATA", reflect.TypeOf(data), data)
	
	json_bytes, _ := json.Marshal(data)
	
	//json_bytes, err := json.Marshal(data)
	/*if err != nil {
		fmt.Println(err)
	}*/
	response = headers(response)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func headers(response http.ResponseWriter) http.ResponseWriter {
	response.Header().Set("Content-Type", "application/json")
	return response
}
