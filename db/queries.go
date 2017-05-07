package db

import (
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
)


func run(q *types.Query) (interface{}, *terr.Trace) {
	/*
	if len(q.filter > 0)
	
	
	
	tbs, tr := GetFilter(dbstr, tablestr, key, val)
	if tr != nil {
		return tables, tr
	}
	for _, table := range(tbs) {
		t := Table{table}
		tables = append(tables, &t)
	}*/
	fmt.Println(q)
	return q, nil
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