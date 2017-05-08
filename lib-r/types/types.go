package types


type Db struct {
	Name string `json:"name"`
}

type Table struct {
	Name string `json:"name"`
}

type Filter struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Query struct {
	Db string
	Table string
	Filters []Filter
	Limit int
}

type Doc struct {
	Data interface{} `json:"data"`
}

type Conf struct {
	Addr string
	User string
	Pwd string
}