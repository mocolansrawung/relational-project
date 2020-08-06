package repositories

import (
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
)

type Repository struct {
	Db *infras.MysqlConn `inject:"db"`
}

func (r *Repository) Get() {
	fmt.Println("Test From Repo")
}
