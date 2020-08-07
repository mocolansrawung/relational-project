package repositories

import "github.com/evermos/boilerplate-go/infras"

type ExampleRepository struct {
	Db *infras.MysqlConn `inject:"db"`
}

func (r *ExampleRepository) Get() (string, error) {
	var status []string
	if err := r.Db.Read.Select(&status, "SELECT 'good' as status"); err != nil {
		return "", err
	}

	return status[0], nil
}
