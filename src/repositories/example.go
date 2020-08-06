package repositories

import "github.com/evermos/boilerplate-go/infras"

type ExampleRepository struct {
	Db *infras.MysqlConn `inject:"db"`
}

func (r *ExampleRepository) Get() string {
	return "good"
}
