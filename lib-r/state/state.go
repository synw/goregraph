package state

import (
	"github.com/synw/goregraph/lib-r/conf"
	"github.com/synw/goregraph/lib-r/types"
	"github.com/synw/terr"
	"net/http"
)

var Conf *types.Conf
var HttpServer *http.Server
var Verbosity int
var DbType string
var Dbs []string
var Db string
var Table string

func InitState(dev bool, verbosity int, config ...*types.Conf) *terr.Trace {
	Verbosity = verbosity
	// config
	if len(config) == 1 {
		Conf = config[0]
		return nil
	} else {
		cf, tr := conf.GetConf(dev, verbosity)
		if tr != nil {
			return tr
		}
		Conf = cf
	}
	return nil
}
