package main

import (
	"flag"
	"fmt"
	"github.com/synw/goregraph/db"
	"github.com/synw/goregraph/lib-r/httpServer"
	"github.com/synw/goregraph/lib-r/state"
)

var dev = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 0, "Verbosity")

func main() {
	flag.Parse()
	// init state
	state.InitState(*dev, *verbosity)
	// init db
	err := db.Init(state.Conf, false)
	if err != nil {
		fmt.Println(err)
	}
	// run http server
	defer httpServer.Stop()
	if *verbosity > 0 {
		defer fmt.Println("Exit")
	}
	httpServer.InitHttpServer(true)
	select {}
}
