package types


type Db struct {
	Name string `json:"name"`
}

type Table struct {
	Name string `json:"name"`
}

type Filter struct {
	Key string
	Val interface{}
}

type Query struct {
	Db string
	Table string
	Filter
}

