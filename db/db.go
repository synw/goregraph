package db

import (
	"fmt"
	r "gopkg.in/dancannon/gorethink.v3"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
	//"encoding/json"
	//"reflect"
)

var conn *r.Session
var verbose = 0


func Init(config *types.Conf, noinit ...bool) error {
	if len(noinit) == 0 {
		state.InitState(config.Dev, config.Verb, config)
	}
	tr := initDb()
	if tr != nil {
		err := tr.ToErr()
		return err
	}
	dbs, tr := GetDbs()
	if tr != nil {
		err := tr.ToErr()
		return err
	}
	state.Dbs = dbs
	// print message if verbose option
	ready()
	return nil
}

func GetTables(db string) ([]string, *terr.Trace) {
	var tables []string
	res, err := r.DB(db).TableList().Run(conn)
	if err != nil {
		tr := terr.New("db.GetTables", err)
		return tables, tr
	}
	var row interface{}
	for res.Next(&row) {
	    tables = append(tables, row.(string))
	}
	if res.Err() != nil {
	    tr := terr.New("db.GetTables", err)
		return tables, tr
	}
	return tables, nil
}

func GetDbs() ([]string, *terr.Trace) {
	var dbs []string
	res, err := r.DBList().Run(conn)
	if err != nil {
		tr := terr.New("db.GetDbs", err)
		return dbs, tr
	}
	var row interface{}
	for res.Next(&row) {
	    dbs = append(dbs, row.(string))
	}
	if res.Err() != nil {
	    tr := terr.New("db.GetDbs", err)
		return dbs, tr
	}
	return dbs, nil
}

// internal methods

func initDb() *terr.Trace {
	cn, tr := connect()
	if tr != nil {
		tr := terr.Pass("db.InitDb", tr)
		return tr
	}
	conn = cn
	return nil
}

func connect() (*r.Session, *terr.Trace) {
	user := state.Conf.User
	pwd := state.Conf.Pwd
	addr := state.Conf.Addr
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		Username: user,
		Password: pwd,
		InitialCap: 10,
        MaxOpen:    10,
	})
    if err != nil {
        tr := terr.New("db.rethinkdb.connectToDb()", err)
        terr.Fatal("db.connect", tr)
        return session, tr
    }
    return session, nil
}

func ready() {
	if state.Verbosity > 0 {
		fmt.Println("Database ready at", state.Conf.Addr)
	}
}
