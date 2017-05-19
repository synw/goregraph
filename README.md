# Goregraph

Turn a schemaless document database into an Graphql server in one minute: just connect a database and run the server. 
It is also possible to use this package as a library to run Reql queries from Graphql queries

The goal of this project is to have an API server that can plug on an existing Rethinkdb database and be instantly ready
to serve some basic read only queries from it

## Supported databases

- [x] Rethinkdb
- [ ] Postgresql
- [ ] Sqlite

## Install

   ```bash
   go get github.com/synw/goregraph
   go install github.com/synw/goregraph
   ```

Grab the binary and make a `config.json` file to in the same folder than the binary. Set your database's type, address
and credentials:

   ```javascript
   {
   	"type": "rethinkdb",
	"host": "localhost:8080",
	"addr":"localhost:28015",
	"user":"",
	"password":"",
	"cors": ["*"]
	}
   ```

The `cors` parameter is a list of authorized domains that will receive cors headers in the http responses. `host` is the
http host adress and `addr` is the database location

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
    "log"
    "net/http"
    "github.com/synw/goregraph/lib-r/types"
    "github.com/synw/goregraph/db"
    grg "github.com/synw/goregraph/lib-r/httpServer"
    
   )

   func main() {
    //normal stuff
    http.HandleFunc("/*", MyPageHandler)
    // map your graphql endpoint here
    http.HandleFunc("/graphql", grg.HandleQuery)
    
    // database config
    addr := "localhost:28015"
	user := "admin"
	pwd := "adminpasswd"
	cors := []string{"localhost"}
	verbosity = 0
	conf := &types.Conf{addr, user, pwd, dev, verbosity, cors}
	
    // init and check the database connection
	err := db.Init(conf)
	if err != nil {
		fmt.Println(err)
	}
    
    // done
    log.Fatal(http.ListenAndServe(":8080", nil))
}

   ```

## Available queries:

   ```bash
   # get a list of databases
   curl -g 'http://localhost:8080/graphql?query={dbs{name}}'
   
   # get a list of tables in a database
   curl -g 'http://localhost:8080/graphql?query={tables(db:"rethinkdb"){name}}'
   
   # count documents in a table
   curl -g 'http://localhost:8080/graphql?query={count(db:"rethinkdb",table:"logs"){num}}'
   
   # get all documents from a table
   curl -g 'http://localhost:8080/graphql?query={docs(db:"rethinkdb",table:"server_status"){data}}'
   
   # limit the number of documents to return
   curl -g 'http://localhost:8080/graphql?query={docs(db:"rethinkdb",table:"logs",limit:10){data}}'
   
   # pluck: limit the fields to return
   curl -g 'http://localhost:8080/graphql?query={docs(db:"rethinkdb",table:"server_status", \
   pluck:"name,network,time_connected"){data}}'
   ```

You can use multiple options together like pluck with limit.

Note: the `data` received is a string: you will have to parse it to turn it into json data

## Todo

- [ ] Add options for the http server
- [x] Add cors headers option
- [ ] Add options to limit the dbs and tables that can be queried
- [ ] Better error handling
- [ ] Ratelimit requests
- [ ] Custom schema injection mechanism
- [ ] Consider adding some authentication or token mechanism
- [ ] More queries and query options

## Credits

- [Gorethink](https://github.com/GoRethink/gorethink): Rethinkdb drivers
- [Graphql-go](https://github.com/graphql-go/graphql): Graphql drivers
- [Chi](https://github.com/pressly/chi): http router
- [Cors](https://github.com/goware/cors): cors http headers
- [Terr](https://github.com/synw/terr): error handling

