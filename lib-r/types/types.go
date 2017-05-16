package types


type Db struct {
	Name string `json:"name"`
}

type Table struct {
	Name string `json:"name"`
}

type Query struct {
	Db string
	Table string
	Limit int
	Pluck []string
}

type Doc struct {
	Data string `json:"data"`
}

type Conf struct {
	Addr string
	User string
	Pwd string
	Dev bool
	Verb int
	Cors []string
}