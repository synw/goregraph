# Goregraph

Turn a Rethinkdb database into an Graphql server in one minute

## Install

   ```bash
   go get github.com/synw/goregraph
   go install github.com/synw/goregraph
   ```

Grab the binary and make a `config.json` file to in the same folder than the binary. Set your Rethinkdb database's address
and credentials:

   ```javascript
   {
	"addr":"localhost:28015",
	"user":"",
	"password":""
	}
   ```

## Run the Graphql server

   ```bash
   ./goregraph
   ```

The server is ready for queries at `http://localhost:8080`

Check the [available queries](https://github.com/synw/goregraph#available-queries)

## Use as library

   ```go
   package main
   
   import (
      "github.com/synw/goregraph/db"
      g "github.com/synw/goregraph/lib-r/types"
   )
   
   
   func main() {
      // ...
      // initialize the database connection
      config := &g.Conf("localhost:28015", "db_user", "db_password"}
      err := db.Init(config)
      if err != nil {
         fmt.Println(err)
      }
      // ...
   }
   
   func handleQuery(response http.ResponseWriter, request *http.Request) {
	  q := request.URL.Query()["query"][0]
	  result, tr := db.RunQuery(q)
	  if tr != nil {
	     // tr is a custom stack trace from goregraph
	     fmt.Println(tr.Print())
	  	 // to translate it to an error your can do: tr.ToErr()
	  }
	  json_bytes, err := json.Marshal(result.Data)
	  if err != nil {
	     fmt.Println(err)
	  }
	  response.Header().Set("Content-Type", "application/json")
	  fmt.Fprintf(response, "%s\n", json_bytes)
   }
   ```

## Available queries:

   ```bash
   # get a list of databases
   curl -g 'http://localhost:8080/graphql?query={dbs{name}}'
   
   # get a list of tables in a database
   curl -g 'http://localhost:8080/graphql?query={tables(db:"rethinkdb"){name}}'
   
   # get some documents
   curl -g 'http://localhost:8080/graphql?query={getAll(db:"rethinkdb",table:"logs",limit:20){data,id}}' 
   ```
   
## Todo

- [ ] Add options for the http server
- [ ] Add cors headers option
- [ ] Add options to limit the dbs and tables that can be queried
- [ ] Ratelimit requests
- [ ] More queries
