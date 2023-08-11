package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(product Product) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (product Product, err error)
	Update(product Product) (err error)
	ResolveImagesByProductIDs(ids []uuid.UUID) (images []ProductImage, err error)
	ResolveVariantsByProductIDs(ids []uuid.UUID) (variants []ProductVariant, err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	s := new(ProductRepositoryMySQL)
	s.DB = db
	return s
}

func (r *ProductRepositoryMySQL) Create(product Product) (err error) {
	exists, err := r.ExistsByID(product.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "product", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, product); err != nil {
			e <- err
			return
		}

		if err := r.txCreateImages(tx, product.Images); err != nil {
			e <- err
			return
		}

		if err := r.txCreateVariants(tx, product.Variants); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (r *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM products WHERE id = ?",
		id.String())

	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (r *ProductRepositoryMySQL) ResolveByID(id uuid.UUID) (product Product, err error) {
	query := `SELECT * FROM products`
	err = r.DB.Read.Get(
		&product,
		query+" WHERE id = ?",
		id.String())

	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("product")
		logger.ErrorWithStack(err)
		return
	}

	return
}
func (r *ProductRepositoryMySQL) Update(product Product) (err error) {
	exists, err := r.ExistsByID(product.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if !exists {
		err = failure.NotFound("product")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txDeleteImages(tx, product.ID); err != nil {
			e <- err
			return
		}
		if err := r.txDeleteVariants(tx, product.ID); err != nil {
			e <- err
			return
		}

		if err := r.txCreateImages(tx, product.Images); err != nil {
			e <- err
			return
		}
		if err := r.txCreateVariants(tx, product.Variants); err != nil {
			e <- err
			return
		}

		if err := r.txUpdate(tx, product); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (r *ProductRepositoryMySQL) ResolveImagesByProductIDs(ids []uuid.UUID) (images []ProductImage, err error) {
	if len(ids) == 0 {
		return
	}

	imageQuery := `SELECT * FROM images`

	query, args, err := sqlx.In(imageQuery+" WHERE product_id IN (?)", ids)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	err = r.DB.Read.Select(&images, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
func (r *ProductRepositoryMySQL) ResolveVariantsByProductIDs(ids []uuid.UUID) (variants []ProductVariant, err error) {
	if len(ids) == 0 {
		return
	}

	imageQuery := `SELECT * FROM variants`

	query, args, err := sqlx.In(imageQuery+" WHERE product_id IN (?)", ids)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	err = r.DB.Read.Select(&variants, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

// Internal Functions
func (r *ProductRepositoryMySQL) txCreate(tx *sqlx.Tx, product Product) (err error) {
	query := `INSERT INTO products (id, brand_id, user_id, name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	) VALUES ( :id, :brand_id, :user_id, :name, :created_at, :created_by, :updated_at, :updated_by, :deleted_at, :deleted_by)`

	stmt, err := tx.PrepareNamed(query)
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
func (r *ProductRepositoryMySQL) composeBulkInsertImageQuery(images []ProductImage) (query string, params []interface{}, err error) {
	queryBulk := `INSERT INTO images (id, product_id, image_url, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) VALUES `
	queryBulkPlaceholder := `(:id, :product_id, :image_url, :created_at, :created_by, :updated_at, :updated_by, :deleted_at, :deleted_by)`

	values := []string{}
	for _, i := range images {
		param := map[string]interface{}{
			"id":         i.ID,
			"product_id": i.ProductID,
			"image_url":  i.ImageURL,
			"created_at": i.CreatedAt,
			"created_by": i.CreatedBy,
			"updated_at": i.UpdatedAt,
			"updated_by": i.UpdatedBy,
			"deleted_at": i.DeletedAt,
			"deleted_by": i.DeletedBy,
		}

		q, args, err := sqlx.Named(queryBulk, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}

	query = fmt.Sprintf("%v %v", queryBulkPlaceholder, strings.Join(values, ","))

	return
}
func (r *ProductRepositoryMySQL) txCreateImages(tx *sqlx.Tx, images []ProductImage) (err error) {
	if len(images) == 0 {
		return
	}

	query, args, err := r.composeBulkInsertImageQuery(images)
	if err != nil {
		return
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (r *ProductRepositoryMySQL) composeBulkInsertVariantQuery(variants []ProductVariant) (query string, params []interface{}, err error) {
	queryBulk := `INSERT INTO variants (id, product_id, name, price, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) VALUES `
	queryBulkPlaceholder := `(:id, :product_id, :name, :price, :created_at, :created_by, :updated_at, :updated_by, :deleted_at, :deleted_by)`

	values := []string{}
	for _, v := range variants {
		param := map[string]interface{}{
			"id":         v.ID,
			"product_id": v.ProductID,
			"name":       v.Name,
			"price":      v.Price,
			"created_at": v.CreatedAt,
			"created_by": v.CreatedBy,
			"updated_at": v.UpdatedAt,
			"updated_by": v.UpdatedBy,
			"deleted_at": v.DeletedAt,
			"deleted_by": v.DeletedBy,
		}

		q, args, err := sqlx.Named(queryBulk, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}

	query = fmt.Sprintf("%v %v", queryBulkPlaceholder, strings.Join(values, ","))

	return
}
func (r *ProductRepositoryMySQL) txCreateVariants(tx *sqlx.Tx, variants []ProductVariant) (err error) {
	if len(variants) == 0 {
		return
	}

	query, args, err := r.composeBulkInsertVariantQuery(variants)
	if err != nil {
		return
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (r *ProductRepositoryMySQL) txDeleteImages(tx *sqlx.Tx, productID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM images WHERE product_id = ?", productID.String())

	return
}
func (r *ProductRepositoryMySQL) txDeleteVariants(tx *sqlx.Tx, productID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM variants WHERE product_id = ?", productID.String())

	return
}
func (r *ProductRepositoryMySQL) txUpdate(tx *sqlx.Tx, product Product) (err error) {
	query := `UPDATE products SET id = :id, brand_id = :brand_id, user_id = :user_id, name = :name, created_at = :created_at, created_by = :created_by, updated_at = :updated_at, updated_by = :updated_by, deleted_at = :deleted_at, deleted_by = :deleted_by`
	stmt, err := tx.PrepareNamed(query)
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
