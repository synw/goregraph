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

func InitState(name string, verbosity int, config ...*types.Conf) *terr.Trace {
	Verbosity = verbosity
	// config
	if len(config) == 1 {
		Conf = config[0]
		return nil
	}
	cf, tr := conf.GetConf(name, verbosity)
	if tr != nil {
		return tr
	}
	Conf = cf
	return nil
}