package rethinkdb

import (
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/terr"
	r "gopkg.in/dancannon/gorethink.v3"
)

var conn *r.Session

func InitDb() *terr.Trace {
	cn, tr := connect()
	if tr != nil {
		tr := terr.Pass("db.rethinkdb.InitDb", tr)
		return tr
	}
	conn = cn
	return nil
}

func GetDbs() ([]string, *terr.Trace) {
	var dbs []string
	res, err := r.DBList().Run(conn)
	if err != nil {
		tr := terr.New("db.rethinkdb.GetDbs", err)
		return dbs, tr
	}
	var row interface{}
	for res.Next(&row) {
		dbs = append(dbs, row.(string))
	}
	if res.Err() != nil {
		tr := terr.New("db.rethinkdb.GetDbs", err)
		return dbs, tr
	}
	return dbs, nil
}

func connect() (*r.Session, *terr.Trace) {
	user := state.Conf.User
	pwd := state.Conf.Pwd
	addr := state.Conf.Addr
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address:    addr,
		Username:   user,
		Password:   pwd,
		InitialCap: 10,
		MaxOpen:    10,
	})
	if err != nil {
		tr := terr.New("db.rethinkdb.connect()", err)
		return session, tr
	}
	return session, nil
}
