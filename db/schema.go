package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/synw/goregraph/lib-r/types"
	"github.com/synw/terr"
	"reflect"
	"strings"
)

var dbType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Db",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	},
)

var tableType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Table",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	},
)

var docType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Doc",
		Fields: graphql.Fields{
			"data": &graphql.Field{Type: graphql.String},
		},
	},
)

var countType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Count",
		Fields: graphql.Fields{
			"num": &graphql.Field{Type: graphql.Int},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"dbs": &graphql.Field{
				Type: graphql.NewList(dbType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					dbs := getDbs()
					return dbs, nil
				},
			},
			"tables": &graphql.Field{
				Type: graphql.NewList(tableType),
				Args: graphql.FieldConfigArgument{
					"db": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tables, tr := getTables(p.Args["db"].(string))
					if tr != nil {
						return tables, tr.ToErr()
					}
					return tables, nil
				},
			},
			"count": &graphql.Field{
				Type: countType,
				Args: graphql.FieldConfigArgument{
					"db": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"table": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Args["db"].(string)
					table := p.Args["table"].(string)
					q := &types.CountQuery{db, table, 0}
					num, tr := countDocs(q)
					if tr != nil {
						return num, tr.ToErr()
					}
					return num, nil
				},
			},
			"doc": &graphql.Field{
				Type: docType,
				Args: graphql.FieldConfigArgument{
					"db": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"table": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Args["db"].(string)
					table := p.Args["table"].(string)
					var pluck []string
					q := &types.Query{db, table, 0, pluck}
					doc, tr := getDoc(q)
					if tr != nil {
						return doc.Data, tr.ToErr()
					}
					fmt.Println("DOC", reflect.TypeOf(doc), doc.Data)
					return doc, nil
				},
			},
			"docs": &graphql.Field{
				Type: graphql.NewList(docType),
				Args: graphql.FieldConfigArgument{
					"db": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"table": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"pluck": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Args["db"].(string)
					table := p.Args["table"].(string)
					limitArg := p.Args["limit"]
					limit := 0
					if limitArg != nil {
						limit = limitArg.(int)
					}
					pluckArg := p.Args["pluck"]
					var pluck []string
					if pluckArg != nil {
						pluckStr := pluckArg.(string)
						i := strings.Index(pluckStr, ",")
						if i > -1 {
							pluck = strings.Split(pluckStr, ",")
						} else {
							pluck = []string{pluckStr}
						}
					}
					q := &types.Query{db, table, limit, pluck}
					res, tr := getDocs(q)
					if tr != nil {
						return res, tr.ToErr()
					}
					var data []*types.Doc
					for _, doc := range res {
						data = append(data, doc)
						//fmt.Println("ELEM", doc.Data[0:15])
					}
					return data, nil
				},
			},
		},
	})

func RunQuery(q string) ([]byte, *terr.Trace) {
	res := graphql.Do(graphql.Params{
		Schema:        Schem,
		RequestString: q,
	})
	b := []byte{}
	if len(res.Errors) > 0 {
		msg := fmt.Sprintf("wrong res, unexpected errors: %v", res.Errors)
		err := errors.New(msg)
		tr := terr.New("db.schema.RunQuery", err)
		tr.Printf("db.schema.RunQuery")
		return b, tr
	}
	data := res.Data
	json_bytes, _ := json.Marshal(data)
	//fmt.Println(res.Data)

	return json_bytes, nil
}

var Schem, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)
