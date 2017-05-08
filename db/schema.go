package db

import (
	//"fmt"
	"github.com/graphql-go/graphql"
	"github.com/synw/goregraph/lib-r/types"
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

var getAllType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Filter",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.String},
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
			"getAll": &graphql.Field{
				Type: graphql.NewList(getAllType),
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
					limit := p.Args["limit"].(int)
					filter := types.Filter{}
					filters := []types.Filter{filter}
					q := &types.Query{db,  table,  filters, limit}
					res, tr := run(q)
					if tr != nil {
						return res, tr.ToErr()
					}
					return res, nil
				},
			},
		},
	})

var Schem, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)
