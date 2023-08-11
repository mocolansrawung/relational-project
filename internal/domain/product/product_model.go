package product

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// Product
type Product struct {
	ID        uuid.UUID        `db:"id" validate:"required"`
	BrandID   uuid.UUID        `db:"brand_id" validate:"required"`
	UserID    uuid.UUID        `db:"user_id" validate:"required"`
	Name      string           `db:"name" validate:"required"`
	Brand     string           `db:"-" validate:"required"`
	Stock     int              `db:"-" validate:"required"`
	CreatedAt time.Time        `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID        `db:"created_by" validate:"required"`
	UpdatedAt null.Time        `db:"updated_at"`
	UpdatedBy nuuid.NUUID      `db:"updated_by"`
	DeletedAt null.Time        `db:"deleted_at"`
	DeletedBy nuuid.NUUID      `db:"deleted_by"`
	Images    []ProductImage   `db:"-" validate:"required,dive,required"`
	Variants  []ProductVariant `db:"-" validate:"required,dive,required"`
}

// Attach Images and Variants to This Product
func (p *Product) AttachImages(images []ProductImage) Product {
	for _, image := range images {
		if image.ProductID == p.ID {
			p.Images = append(p.Images, image)
		}
	}

	return *p
}

func (p *Product) AttachVariants(variants []ProductVariant) Product {
	for _, variant := range variants {
		if variant.ProductID == p.ID {
			p.Variants = append(p.Variants, variant)
		}
	}

	return *p
}

func (p *Product) IsDeleted() (deleted bool) {
	return p.DeletedAt.Valid && p.DeletedBy.Valid
}

func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}

func (p Product) NewFromRequestFormat(req ProductRequestFormat, userID uuid.UUID) (newProduct Product, err error) {
	productID, err := uuid.NewV4()
	if err != nil {
		return newProduct, failure.InternalError(err)
	}

	newProduct = Product{
		ID:        productID,
		Name:      req.Name,
		Brand:     req.Brand,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	images := make([]ProductImage, 0)
	for _, requestImage := range req.Images {
		image := ProductImage{}
		image, _ = image.NewFromRequestFormat(requestImage, productID, userID)
		images = append(images, image)
	}

	variants := make([]ProductVariant, 0)
	for _, requestVariant := range req.Variants {
		variant := ProductVariant{}
		variant, _ = variant.NewFromRequestFormat(requestVariant, productID)
		variants = append(variants, variant)
	}

	newProduct.Images = images
	newProduct.Variants = variants

	err = newProduct.Validate()

	return
}

func (p *Product) Recalculate() {
	p.Stock = int(0)
	recalculatedVariants := make([]ProductVariant, 0)
	for _, variant := range p.Variants {
		variant.Recalculate()
		recalculatedVariants = append(recalculatedVariants, variant)
		p.Stock += variant.Stock
	}

	p.Variants = recalculatedVariants
}

func (p *Product) SoftDelete(userID uuid.UUID) (err error) {
	if p.IsDeleted() {
		return failure.Conflict("softDelete", "product", "already marked as deleted")
	}

	p.DeletedAt = null.TimeFrom(time.Now())
	p.DeletedBy = nuuid.From(userID)

	return
}

func (p Product) ToResponseFormat() ProductResponseFormat {
	resp := ProductResponseFormat{
		ID:        p.ID,
		Name:      p.Name,
		Brand:     p.Brand,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy.Ptr(),
		DeletedAt: p.DeletedAt,
		DeletedBy: p.DeletedBy.Ptr(),
		Images:    make([]ProductImageResponseFormat, 0),
		Variants:  make([]ProductVariantResponseFormat, 0),
	}

	for _, image := range p.Images {
		resp.Images = append(resp.Images, image.ToResponseFormat())
	}

	for _, variant := range p.Variants {
		resp.Variants = append(resp.Variants, variant.ToResponseFormat())
	}

	return resp
}

func (p *Product) Update(req ProductRequestFormat, userID uuid.UUID) (err error) {
	images := make([]ProductImage, 0)
	for _, requestImage := range req.Images {
		image := ProductImage{}
		image, _ = image.NewFromRequestFormat(requestImage, p.ID, p.UserID)
		images = append(images, image)
	}

	variants := make([]ProductVariant, 0)
	for _, requestVariant := range req.Variants {
		variant := ProductVariant{}
		variant, _ = variant.NewFromRequestFormat(requestVariant, p.UserID)
		variants = append(variants, variant)
	}

	p.Images = images
	p.Variants = variants
	p.UpdatedAt = null.TimeFrom(time.Now())
	p.UpdatedBy = nuuid.From(userID)

	p.Recalculate()
	err = p.Validate()

	return
}

func (p *Product) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(p)
}

type ProductRequestFormat struct {
	Name     string                        `json:"name" validate:"required"`
	Brand    string                        `json:"brand" validate:"required"`
	Images   []ProductImageRequestFormat   `json:"images" validate:"required,dive,required"`
	Variants []ProductVariantRequestFormat `json:"variants" validate:"required,dive,required"`
}

type ProductResponseFormat struct {
	ID        uuid.UUID                      `json:"id"`
	BrandID   uuid.UUID                      `json:"-"`
	UserID    uuid.UUID                      `json:"-"`
	Name      string                         `json:"name"`
	Brand     string                         `json:"brand"`
	Stock     int                            `json:"stock"`
	CreatedAt time.Time                      `json:"created_at"`
	CreatedBy uuid.UUID                      `json:"created_by"`
	UpdatedAt null.Time                      `json:"updated_at"`
	UpdatedBy *uuid.UUID                     `json:"updated_by"`
	DeletedAt null.Time                      `json:"deleted_at,omitempty"`
	DeletedBy *uuid.UUID                     `json:"deleted_by,omitempty"`
	Images    []ProductImageResponseFormat   `json:"images"`
	Variants  []ProductVariantResponseFormat `json:"variants"`
}

// ProductImage
type ProductImage struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	ProductID uuid.UUID   `db:"product_id" validate:"required"`
	ImageURL  string      `db:"image_url" validate:"required"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

func (pi ProductImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(pi.ToResponseFormat())
}

func (pi ProductImage) NewFromRequestFormat(format ProductImageRequestFormat, productID uuid.UUID, userID uuid.UUID) (productImage ProductImage, err error) {
	productImageID, err := uuid.NewV4()
	if err != nil {
		return productImage, failure.InternalError(err)
	}

	productImage = ProductImage{
		ID:        productImageID,
		ProductID: productID,
		ImageURL:  format.ImageURL,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	return
}

func (pi *ProductImage) ToResponseFormat() ProductImageResponseFormat {
	return ProductImageResponseFormat{
		ID:        pi.ID,
		ImageURL:  pi.ImageURL,
		CreatedAt: pi.CreatedAt,
		CreatedBy: pi.CreatedBy,
		UpdatedAt: pi.UpdatedAt,
		UpdatedBy: pi.UpdatedBy.Ptr(),
		DeletedAt: pi.DeletedAt,
		DeletedBy: pi.DeletedBy.Ptr(),
	}
}

type ProductImageRequestFormat struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	ImageURL string    `json:"image_url"`
}

type ProductImageResponseFormat struct {
	ID        uuid.UUID  `json:"id"`
	ImageURL  string     `json:"image_url"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy uuid.UUID  `json:"created_by"`
	UpdatedAt null.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	DeletedAt null.Time  `json:"deleted_at"`
	DeletedBy *uuid.UUID `json:"deleted_by"`
}

// ProductVariant
type ProductVariant struct {
	ID                uuid.UUID          `db:"id" validate:"required"`
	ProductID         uuid.UUID          `db:"product_id" validate:"required"`
	Name              string             `db:"name" validate:"required"`
	Stock             int                `db:"-" validate:"required"`
	Price             int                `db:"price" validate:"required"`
	CreatedAt         time.Time          `db:"created_at" validate:"required"`
	CreatedBy         uuid.UUID          `db:"created_by" validate:"required"`
	UpdatedAt         null.Time          `db:"updated_at"`
	UpdatedBy         nuuid.NUUID        `db:"updated_by"`
	DeletedAt         null.Time          `db:"deleted_at"`
	DeletedBy         nuuid.NUUID        `db:"deleted_by"`
	VariantWarehouses []VariantWarehouse `db:"-" validate:"required,dive,required"`
}

func (pv *ProductVariant) AttachVariantWarehouses(VariantWarehouses []VariantWarehouse) ProductVariant {
	for _, vw := range VariantWarehouses {
		if vw.VariantID == pv.ID {
			pv.VariantWarehouses = append(pv.VariantWarehouses, vw)
		}
	}

	return *pv
}

func (pv *ProductVariant) IsDeleted() (deleted bool) {
	return pv.DeletedAt.Valid && pv.DeletedBy.Valid
}

func (pv ProductVariant) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.ToResponseFormat())
}

func (pv ProductVariant) NewFromRequestFormat(req ProductVariantRequestFormat, userID uuid.UUID) (newProductVariant ProductVariant, err error) {
	variantID, err := uuid.NewV4()
	if err != nil {
		return newProductVariant, failure.InternalError(err)
	}

	newProductVariant = ProductVariant{
		ID:        variantID,
		ProductID: req.ProductID,
		Name:      req.Name,
		Price:     req.Price,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	variantWarehouses := make([]VariantWarehouse, 0)
	for _, requestVariantWarehouse := range req.VariantWarehouses {
		variantWarehouse := VariantWarehouse{}
		variantWarehouse, err = variantWarehouse.NewFromRequestFormat(requestVariantWarehouse, variantID, userID)
		if err != nil {
			return
		}

		variantWarehouses = append(variantWarehouses, variantWarehouse)
	}

	newProductVariant.VariantWarehouses = variantWarehouses

	err = newProductVariant.Validate()

	return
}

func (pv *ProductVariant) Recalculate() {
	pv.Stock = int(0)
	recalculatedVariantWarehouses := make([]VariantWarehouse, 0)
	for _, vw := range pv.VariantWarehouses {
		recalculatedVariantWarehouses = append(recalculatedVariantWarehouses, vw)
		pv.Stock += vw.Stock
	}

	pv.VariantWarehouses = recalculatedVariantWarehouses
}

func (pv ProductVariant) ToResponseFormat() ProductVariantResponseFormat {
	resp := ProductVariantResponseFormat{
		ID:                pv.ID,
		Name:              pv.Name,
		Stock:             pv.Stock,
		Price:             pv.Price,
		CreatedAt:         pv.CreatedAt,
		CreatedBy:         pv.CreatedBy,
		UpdatedAt:         pv.UpdatedAt,
		UpdatedBy:         pv.UpdatedBy.Ptr(),
		DeletedAt:         pv.DeletedAt,
		DeletedBy:         pv.DeletedBy.Ptr(),
		VariantWarehouses: make([]VariantWarehouseResponseFormat, 0),
	}

	for _, vw := range pv.VariantWarehouses {
		resp.VariantWarehouses = append(resp.VariantWarehouses, vw.ToResponseFormat())
	}

	return resp
}

func (pv *ProductVariant) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(pv)
}

type ProductVariantRequestFormat struct {
	ProductID         uuid.UUID                       `json:"product_id" validate:"required"`
	Name              string                          `json:"name" validate:"required"`
	Price             int                             `json:"price" validate:"required,min=0"`
	VariantWarehouses []VariantWarehouseRequestFormat `json:"variant_warehouse" validate:"required,dive,required"`
}

type ProductVariantResponseFormat struct {
	ID                uuid.UUID                        `json:"id"`
	Name              string                           `json:"name"`
	Stock             int                              `json:"stock"`
	Price             int                              `json:"price"`
	CreatedAt         time.Time                        `json:"created_at"`
	CreatedBy         uuid.UUID                        `json:"created_by"`
	UpdatedAt         null.Time                        `json:"updated_at"`
	UpdatedBy         *uuid.UUID                       `json:"updated_by"`
	DeletedAt         null.Time                        `json:"deleted_at"`
	DeletedBy         *uuid.UUID                       `json:"deleted_by"`
	VariantWarehouses []VariantWarehouseResponseFormat `json:"variant_warehouse"`
}

// Variant Warehouse Stock
type VariantWarehouse struct {
	VariantID   uuid.UUID   `db:"variant_id" validate:"required"`
	WarehouseID uuid.UUID   `db:"warehouse_id" validate:"required"`
	Stock       int         `db:"stock" validate:"required"`
	CreatedAt   time.Time   `db:"created_at" validate:"required"`
	CreatedBy   uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt   null.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
}

func (vw *VariantWarehouse) IsDeleted() (deleted bool) {
	return vw.DeletedAt.Valid && vw.DeletedBy.Valid
}

func (vw VariantWarehouse) MarshalJSON() ([]byte, error) {
	return json.Marshal(vw.ToResponseFormat())
}

func (vw VariantWarehouse) NewFromRequestFormat(format VariantWarehouseRequestFormat, variantID uuid.UUID, userID uuid.UUID) (variantWarehouse VariantWarehouse, err error) {
	warehouseID, err := uuid.NewV4()
	if err != nil {
		return variantWarehouse, failure.InternalError(err)
	}

	variantWarehouse = VariantWarehouse{
		VariantID:   variantID,
		WarehouseID: warehouseID,
		Stock:       format.Stock,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	return
}

func (vw *VariantWarehouse) ToResponseFormat() VariantWarehouseResponseFormat {
	return VariantWarehouseResponseFormat{
		VariantID:   vw.VariantID,
		WarehouseID: vw.WarehouseID,
		Stock:       vw.Stock,
		CreatedAt:   vw.CreatedAt,
		CreatedBy:   vw.CreatedBy,
		UpdatedAt:   vw.UpdatedAt,
		UpdatedBy:   vw.UpdatedBy.Ptr(),
		DeletedAt:   vw.DeletedAt,
		DeletedBy:   vw.DeletedBy.Ptr(),
	}
}

type VariantWarehouseRequestFormat struct {
	Stock int `json:"stock" validate:"required"`
}

type VariantWarehouseResponseFormat struct {
	VariantID   uuid.UUID  `json:"variant_id"`
	WarehouseID uuid.UUID  `json:"warehouse_id"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedAt   null.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by"`
	DeletedAt   null.Time  `json:"deleted_at"`
	DeletedBy   *uuid.UUID `json:"deleted_by"`
}

// Supporting Entities (delete later)
type User struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	Type      string      `db:"type" validate:"required,eq=ADMIN|eq=REGULAR"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
	Products  []Product   `db:"-" validate:"required,dive,required"`
}

type Brand struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	Name      string      `db:"brand_name" validate:"required"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at" validate:"required"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
	Products  []Product   `db:"-" validate:"required,dive,required"`
}

type Warehouse struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	Name      string      `db:"name" validate:"required"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}
