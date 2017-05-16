package db

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/synw/goregraph/lib-r/types"
	"reflect"
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
					//id := p.Args["id"].(string)		
					filter := types.Filter{}
					filters := []types.Filter{filter}
					q := &types.Query{db,  table,  filters, 0}
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
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := p.Args["db"].(string)
					table := p.Args["table"].(string)
					limitI := p.Args["limit"]
					limit := 0
					if limitI != nil {
						limit = limitI.(int)
					}
					filter := types.Filter{}
					filters := []types.Filter{filter}
					q := &types.Query{db,  table,  filters, limit}
					res, tr := getDocs(q)
					if tr != nil {
						return res, tr.ToErr()
					}
					var data []*types.Doc
					for _, doc := range(res) {
						data = append(data, doc)
						//fmt.Println("ELEM", doc.Data[0:15])
					}
					return data, nil
				},
			},
		},
	})

var Schem, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)
