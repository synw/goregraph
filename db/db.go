package db

import (
	"fmt"
	"errors"
	r "gopkg.in/dancannon/gorethink.v3"
	"github.com/graphql-go/graphql"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
)

var conn *r.Session
var verbose = 0

func RunQuery(query string) (*graphql.Result, *terr.Trace) {
	result := graphql.Do(graphql.Params{
		Schema: Schem,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		msg := fmt.Sprintf("wrong result, unexpected errors: %v", result.Errors)
		err := errors.New(msg)
		tr := terr.New("schema.ExecuteQuery", err)
		return result, tr
	}
	return result, nil
}

func InitVerbose(dev_mode ...string) error {
	verbose = 1
	dm := "normal"
	if len(dev_mode) > 0 {
		dm = dev_mode[0]
	}
	err := Init(dm)
	if err != nil {
		return err
	}
	fmt.Println("Database ready at", state.Addr)
	return nil
}

func Init(dev_mode ...string) error {
	dm := "normal"
	if len(dev_mode) > 0 {
		dm = dev_mode[0]
	}
	// init state
	tr := state.InitState(dm, verbose)
	if tr != nil {
		fmt.Println(tr.Formatc())
	}
	// init db
	tr = initDb()
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
	user := state.User
	pwd := state.Pwd
	addr := state.Addr
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
        return session, tr
    }
    return session, nil
}
