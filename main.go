package main

import (
	"fmt"
	"flag"
	"github.com/synw/goregraph/lib-r/httpServer"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/db"
)


var dev = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 0, "Verbosity")

func main() {
	flag.Parse()
	// init state
	state.InitState(*dev, *verbosity)
	// init db
	db.Init(state.Conf, false)
	// run http server
	defer httpServer.Stop() 
	if *verbosity > 0 { 
		defer fmt.Println("Exit") 
	}
	httpServer.InitHttpServer(true)
}
