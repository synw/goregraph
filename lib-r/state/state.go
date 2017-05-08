package state

import (
	"net/http"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/conf"
	"github.com/synw/goregraph/lib-r/types"
)


var Conf *types.Conf
var HttpServer *http.Server
var Verbosity int
var Dbs []string
var Db string
var Table string

func InitState(name string, verbosity int) *terr.Trace {
	Verbosity = verbosity
	// config
	cf, tr := conf.GetConf(name)
	if tr != nil {
		return tr
	}
	Conf = cf
	return nil
}