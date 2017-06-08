package rethinkdb

import (
	"fmt"
	//"strings"
	"github.com/Jeffail/gabs"
	"github.com/synw/goregraph/lib-r/types"
	"github.com/synw/terr"
	r "gopkg.in/dancannon/gorethink.v3"
)

func CountDocs(q *types.CountQuery) (*types.Count, *terr.Trace) {
	count := &types.Count{}
	reql := r.DB(q.Db).Table(q.Table).Count()
	res, err := reql.Run(conn)
	if err != nil {
		tr := terr.New("db.countDocs", err)
		return count, tr
	}
	var row interface{}
	res.One(&row)
	count.Data = int(row.(float64))
	return count, nil
}

func GetDoc(q *types.Query) (*types.Doc, *terr.Trace) {
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

func GetDocs(q *types.Query) ([]*types.Doc, *terr.Trace) {
	var docs []*types.Doc
	reql := r.DB(q.Db).Table(q.Table)
	if q.Limit > 0 {
		reql = reql.Limit(q.Limit)
	}
	if len(q.Pluck) > 0 {
		reql = reql.Pluck(q.Pluck)
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

func GetTables(db string) ([]types.Table, *terr.Trace) {
	var tables []types.Table
	res, err := r.DB(db).TableList().Run(conn)
	if err != nil {
		tr := terr.New("db.rethinkdb.GetTables", err)
		return tables, tr
	}
	var row interface{}
	for res.Next(&row) {
		t := row.(string)
		table := types.Table{t}
		tables = append(tables, table)
	}
	if res.Err() != nil {
		tr := terr.New("db.rethinkdb.GetTables", err)
		return tables, tr
	}
	return tables, nil
}
