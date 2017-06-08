package types

type Db struct {
	Name string `json:"name"`
}

type Table struct {
	Name string `json:"name"`
}

type Query struct {
	Db    string
	Table string
	Limit int
	Pluck []string
}

type CountQuery struct {
	Db    string
	Table string
	Num   int
}

type Doc struct {
	Data string `json:"data"`
}

type Count struct {
	Data int `json:"num"`
}

type Conf struct {
	DbType string
	Host   string
	Addr   string
	User   string
	Pwd    string
	Dev    bool
	Verb   int
	Cors   []string
}
