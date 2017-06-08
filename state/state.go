package state

import (
	"errors"
	"github.com/synw/goregraph/db"
	grgState "github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
	grgTypes "github.com/synw/goregraph/lib-r/types"
	"github.com/synw/microb-goregraph/conf"
	"github.com/synw/terr"
)

var Verbosity int
var Conf *types.Conf

func InitState(dev bool, verbosity int) *terr.Trace {
	Verbosity = verbosity
	cf, tr := conf.GetConf(dev, verbosity)
	if tr != nil {
		return tr
	}
	Conf = cf
	grgConf := &grgTypes.Conf{
		"rethinkdb",
		Conf.Host,
		Conf.Addr,
		Conf.User,
		Conf.Pwd,
		Conf.Dev,
		Conf.Verb,
		Conf.Cors,
	}
	tr = grgState.InitState(Conf.Dev, Conf.Verb, grgConf)
	if tr != nil {
		err := errors.New("Unable to initialize goregraph state")
		tr = terr.Add("state.InitState", err)
		return tr
	}
	// init goregraph
	err := db.Init(Conf)
	if err != nil {
		tr := terr.New("state.InitState", err)
		return tr
	}
	return nil
}
