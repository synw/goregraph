package db

import (
	//"fmt"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
)


func run(q *types.Query) ([]*types.Doc, *terr.Trace) {
	docs, tr := getDocs(q)
	if tr != nil {
		tr.Printc()
		return docs, tr
	}
	return docs, nil
}

func getTables(dbstr string) ([]*types.Table, *terr.Trace) {
	var tables []*types.Table
	tbs, tr := GetTables(dbstr)
	if tr != nil {
		return tables, tr
	}
	for _, table := range(tbs) {
		t := &types.Table{table}
		tables = append(tables, t)
	}
	return tables, nil
}

func getDbs() []*types.Db {
	var dbs []*types.Db
	for _, db := range(state.Dbs) {
		d := &types.Db{db}
		dbs = append(dbs, d)
	}
	return dbs
}