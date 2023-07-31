package product

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	productQueries = struct {
		insertProduct string
	}{
		insertProduct: `
			INSERT INTO product (
				id,
				name,
				shop_name,
				price,
				created_at
			) VALUES (
				:id,
				:name,
				:shop_name,
				:price,
				:created_at
			)
		`,
	}
)

// ProductRepository is the repository for Product data
type ProductRepository interface {
	Create(product Product) (err error)
}

// ProductRepositoryMySQL is the MYSQL-backed implementation of FooRepository
type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

// ProvideProductRepositoryMySQL is the provider for this repository
func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	s := new(ProductRepositoryMySQL)
	s.DB = db
	return s
}

// Create creates a new Product
// func (r *ProductRepositoryMySQL) Create(product Product) (err error) {
// 	exists, err := r.ExistsByID(product.ID)
// 	if err != nil {
// 		logger.ErrorWithStack(err)
// 		return
// 	}

// 	if exists {
// 		err = failure.Conflict("create", "product", "already exists")
// 		logger.ErrorWithStack(err)
// 		return
// 	}

// 	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
// 		if err := r.txCreate(tx, product); err != nil {
// 			e <- err
// 			return
// 		}

// 		e <- nil
// 	})
// }

func (r *ProductRepositoryMySQL) Create(product Product) (err error) {
	exists, err := r.ExistsByID(product.ID)
}

// Resolve resolves all Products
// func (r *ProductRepositoryMySQL) Resolve(id uuid.UUID) (product Product, err error) {
// 	err = r.DB.Read.Get(

// 	)
// }

func (r *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM product WHERE id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (r *ProductRepositoryMySQL) txCreate(tx *sqlx.Tx, product Product) (err error) {
	stmt, err := tx.PrepareNamed(productQueries.insertProduct)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
