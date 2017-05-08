package db

import (
	"fmt"
	"errors"
	r "gopkg.in/dancannon/gorethink.v3"
	"github.com/graphql-go/graphql"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
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

func getDocs(q *types.Query) ([]*types.Doc, *terr.Trace) {
	var docs []*types.Doc
	var reql r.Term
	if q.Limit > 0 {
		reql = r.DB(q.Db).Table(q.Table).Limit(q.Limit)
	} else {
		reql = r.DB(q.Db).Table(q.Table)
	}
	res, err := reql.Run(conn)
	if err != nil {
		tr := terr.New("db.getDocs", err)
		return docs, tr
	}
	var row interface{}
	for res.Next(&row) {
		doc := &types.Doc{row}
	    docs = append(docs, doc)
	}
	if res.Err() != nil {
	    tr := terr.New("db.getDocs", err)
		return docs, tr
	}
	return docs, nil
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
        return session, tr
    }
    return session, nil
}

func ready() {
	if state.Verbosity > 0 {
		fmt.Println("Database ready at", state.Conf.Addr)
	}
}
