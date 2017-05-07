package main

import (
	"fmt"
	"flag"
	"github.com/synw/goregraph/lib-r/httpServer"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/db"
)


var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 0, "Verbosity")

func main() {
	flag.Parse()
	name := "normal"
	if *dev_mode == true {
		name = "dev"
	}
	// init state
	state.InitState(name, *verbosity)
	if *verbosity > 0 {
		db.InitVerbose(name)
	} else {
		db.Init(name)
	}
	// run http server
	defer httpServer.Stop() 
	if *verbosity > 0 { 
		defer fmt.Println("Exit") 
	}
	httpServer.InitHttpServer(true)
	select{}
}
