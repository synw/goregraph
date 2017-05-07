package state

import (
	"net/http"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/conf"
)


var Addr string
var User string
var Pwd string
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
	// db credentials
	Addr = cf["addr"].(string)
	User = cf["user"].(string)
	Pwd = cf["password"].(string)
	return nil
}