package db

import (
	"errors"
	"fmt"
	"github.com/synw/goregraph/db/rethinkdb"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
	"github.com/synw/terr"
	//"encoding/json"
	//"reflect"
)

var verbose = 0

func Init(config *types.Conf, noinit ...bool) error {
	if len(noinit) == 0 {
		state.InitState(config.Dev, config.Verb, config)
	}
	var tr *terr.Trace
	var dbs []string
	if state.Conf.DbType == "rethinkdb" {
		tr = rethinkdb.InitDb()
		if tr != nil {
			err := tr.ToErr()
			return err
		}
		dbs, tr = rethinkdb.GetDbs()
		if tr != nil {
			err := tr.ToErr()
			return err
		}
	} else {
		err := errors.New("Database type not implemented")
		tr := terr.New("db.Init", err)
		return tr.ToErr()
	}
	state.Dbs = dbs
	// print message if verbose option
	ready()
	return nil
}

func GetDbs() []string {
	return state.Dbs
}

func getDbs() []*types.Db {
	var dbs []*types.Db
	for _, db := range state.Dbs {
		d := &types.Db{db}
		dbs = append(dbs, d)
	}
	return dbs
}

func countDocs(q *types.CountQuery) (*types.Count, *terr.Trace) {
	var count *types.Count
	if state.Conf.DbType == "rethinkdb" {
		n, tr := rethinkdb.CountDocs(q)
		count = n
		if tr != nil {
			tr = terr.Pass("db.CountDocs", tr)
			return count, tr
		}
	} else {
		err := errors.New("Database type not implemented")
		tr := terr.New("db.CountDocs", err)
		return count, tr
	}
	return count, nil
}

func getDoc(q *types.Query) (*types.Doc, *terr.Trace) {
	var doc *types.Doc
	if state.Conf.DbType == "rethinkdb" {
		doc, tr := rethinkdb.GetDoc(q)
		if tr != nil {
			tr = terr.Pass("db.GetDoc", tr)
			return doc, tr
		}
	} else {
		err := errors.New("Database type not implemented")
		tr := terr.New("db.GetDoc", err)
		return doc, tr
	}
	return doc, nil
}

func getDocs(q *types.Query) ([]*types.Doc, *terr.Trace) {
	var docs []*types.Doc
	if state.Conf.DbType == "rethinkdb" {
		d, tr := rethinkdb.GetDocs(q)
		docs = d
		if tr != nil {
			tr = terr.Pass("db.GetDocs", tr)
			return docs, tr
		}
	} else {
		err := errors.New("Database type not implemented")
		tr := terr.New("db.GetDocs", err)
		return docs, tr
	}
	return docs, nil
}

func getTables(dbstr string) ([]types.Table, *terr.Trace) {
	var tables []types.Table
	if state.Conf.DbType == "rethinkdb" {
		t, tr := rethinkdb.GetTables(dbstr)
		tables = t
		if tr != nil {
			tr = terr.Pass("db.GetTables", tr)
			return tables, tr
		}
	} else {
		err := errors.New("Database type not implemented")
		tr := terr.New("db.GetTables", err)
		return tables, tr
	}
	return tables, nil
}

func ready() {
	if state.Verbosity > 0 {
		fmt.Println("Database ready at", state.Conf.Addr)
	}
}
