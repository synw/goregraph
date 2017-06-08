package mutate

import (
	"errors"
	"github.com/synw/microb-goregraph/server"
	"github.com/synw/terr"
)

func Start() *terr.Trace {
	if server.Server.Running == true {
		err := errors.New("Graphql server is already running")
		tr := terr.New("state.mutate.Start", err)
		return tr
	}
	go server.InitServer()
	return nil
}

func Stop() *terr.Trace {
	if server.Server.Running == false {
		err := errors.New("Graphql server is not running")
		tr := terr.New("state.mutate.Stop", err)
		return tr
	}
	tr := server.Stop()
	if tr != nil {
		return tr
	}
	return nil
}
