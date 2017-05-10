package db

import (
	"fmt"
	r "gopkg.in/dancannon/gorethink.v3"
	"github.com/Jeffail/gabs"
	"github.com/synw/terr"
	"github.com/synw/goregraph/lib-r/state"
	"github.com/synw/goregraph/lib-r/types"
)


func getDoc(q *types.Query) (*types.Doc, *terr.Trace) {
	var reql r.Term
	var doc *types.Doc
	reql = r.DB(q.Db).Table(q.Table)
	res, err := reql.Run(conn)
	if err != nil {
		tr := terr.New("db.getDoc", err)
		return doc, tr
	}
	var row interface{}
	res.One(&row)
    obj, err := gabs.Consume(&row)
    if err != nil {
    	fmt.Println("ERR", err)
    }
	doc = &types.Doc{obj.String()}
	return doc, nil
}

func getDocs(q *types.Query) ([]*types.Doc, *terr.Trace) {
	var reql r.Term
	var docs []*types.Doc
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
	objs := gabs.New()
	objs.Array("docs", "array")
	for res.Next(&row) {
	    obj, err := gabs.Consume(&row)
	    if err != nil {
	    	fmt.Println("ERR", err)
	    }
	    doc := &types.Doc{obj.String()}
		objs.ArrayAppend(doc, "docs", "array")
		docs = append(docs, doc)
    	//fmt.Println("JSON", reflect.TypeOf(doc), doc.Data[0:15])
	}
	//fmt.Println("DOCS", docs)
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