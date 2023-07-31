package product

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// Entity model for this product services
type Product struct {
	ProductID   uuid.UUID        `db:"product_id" validate:"required"`
	BrandID     uuid.UUID        `db:"brand_id" validate:"required"`
	UserID      uuid.UUID        `db:"user_id" validate:"required"`
	ProductName string           `db:"product_name"`
	CreatedAt   time.Time        `db:"created_at"`
	CreatedBy   uuid.UUID        `db:"created_by" validate:"required"`
	UpdatedAt   time.Time        `db:"updated_at"`
	UpdatedBy   uuid.UUID        `db:"updated_by" validate:"required"`
	DeletedAt   null.Time        `db:"deleted_at"`
	DeletedBy   nuuid.NUUID      `db:"deleted_by"`
	Images      []ProductImage   `db:"-" validate:"required,dive,required"`
	Variants    []ProductVariant `db:"-" validate:"required,dive,required"`
}

// AttachImages attaches Images to the Product
func (p *Product) AttachImages(images []ProductImage) Product {
	for _, image := range images {
		if image.ProductID == p.ProductID {
			p.Images = append(p.Images, image)
		}
	}

	return *p
}

// AttachVariants attaches Variants to the Product
func (p *Product) AttachVariants(variants []ProductVariant) Product {
	for _, variant := range variants {
		if variant.ProductID == p.ProductID {
			p.Variants = append(p.Variants, variant)
		}
	}

	return *p
}

// Supporting entities - delete later
type User struct {
	UserID    uuid.UUID   `db:"user_id" validate:"required"`
	UserType  string      `db:"user_type" validate:"required,eq=ADMIN|eq=REGULAR"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
	Products  []Product   `db:"-" validate:"required,dive,required"`
}

type Brand struct {
	BrandID   uuid.UUID   `db:"brand_id" validate:"required"`
	BrandName string      `db:"brand_name" validate:"required"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
	Products  []Product   `db:"-" validate:"required,dive,required"`
}

type Warehouse struct {
	WarehouseID   uuid.UUID   `db:"warehouse_id" validate:"required"`
	WarehouseName string      `db:"warehouse_name" validate:"required"`
	CreatedAt     time.Time   `db:"created_at"`
	CreatedBy     uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt     time.Time   `db:"updated_at"`
	UpdatedBy     nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt     null.Time   `db:"deleted_at"`
	DeletedBy     nuuid.NUUID `db:"deleted_by"`
}

func (p *Product) IsDeleted() (deleted bool) {
	return p.DeletedAt.Valid && p.DeletedBy.Valid
}

func (p *Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}

func (p Product) NewFromRequestFormat(req ProductRequestFormat, userID uuid.UUID, brandID uuid.UUID) (newProduct Product, err error) {
	productID, _ := uuid.NewV4()
	newProduct = Product{
		ProductID:   productID,
		BrandID:     brandID,
		UserID:      userID,
		ProductName: req.ProductName,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userID,
	}

	images := make([]ProductImage, 0)
	for _, requestImage := range req.Images {
		image := ProductImage{}
		image = image.NewFromRequestFormat(requestImage, productID)
		images = append(images, image)
	}
	newProduct.Images = images

	variants := make([]ProductVariant, 0)
	for _, requestVariant := range req.Variants {
		variant := ProductVariant{}
		variant = variant.NewFromRequestFormat(requestVariant, productID)
		variants = append(variants, variant)
	}
	newProduct.Variants = variants

	err = newProduct.Validate()

	return
}

func (p Product) ToResponseFormat() ProductResponseFormat {
	resp := ProductResponseFormat{
		ProductID:   p.ProductID,
		ProductName: p.ProductName,
		Brand:       "", // Retrieve BrandName based on BrandID.
		Stock:       0,  // Initialize the total stock to 0.
		Variants:    p.Variants,
		Images:      make([]ProductImage, 0), // Initialize an empty slice for image URLs.
		CreatedAt:   p.CreatedAt,
		CreatedBy:   p.CreatedBy,
		UpdatedAt:   p.UpdatedAt,
		UpdatedBy:   &p.UpdatedBy,
		DeletedAt:   p.DeletedAt,
		DeletedBy:   &p.DeletedBy,
	}

	// Calculate the total stock by summing up the stock from all VariantWarehouse(s).
	for _, variant := range p.Variants {
		for _, variantWarehouse := range variant.VariantWarehouse {
			resp.Stock += variantWarehouse.Stock
		}
	}

	// Populate the image URLs.
	for _, image := range p.Images {
		resp.Images = append(resp.Images, image.ImageURL)
	}

	return resp
}

func (p *Product) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(p)
}

type ProductRequestFormat struct {
	ProductName string           `json:"product_name" validate:"required"`
	Brand       string           `json:"brand" validate:"required"`
	Variants    []ProductVariant `json:"variants" validate:"required,dive,required"`
	Images      []ProductImage   `json:"images" validate:"required,dive,required"`
}

type ProductResponseFormat struct {
	ProductID   uuid.UUID        `json:"product_id"`
	ProductName string           `json:"product_name"`
	Brand       string           `json:"brand"`
	Stock       int              `json:"stock"`
	Variants    []ProductVariant `json:"variants"`
	Images      []ProductImage   `json:"images"`
	CreatedAt   time.Time        `json:"created_at"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	UpdatedAt   time.Time        `json:"updated_at"`
	UpdatedBy   *uuid.UUID       `json:"updated_by"`
	DeletedAt   null.Time        `json:"deleted_at,omitempty"`
	DeletedBy   *nuuid.NUUID     `json:"deleted_by,omitempty"`
}

// ProductVariant
type ProductVariant struct {
	VariantID        uuid.UUID          `db:"variant_id" validate:"required"`
	ProductID        uuid.UUID          `db:"product_id" validate:"required"`
	VariantName      string             `db:"variant_name"`
	Price            int                `db:"price" validate:"required"`
	CreatedAt        time.Time          `db:"created_at"`
	CreatedBy        uuid.UUID          `db:"created_by" validate:"required"`
	UpdatedAt        time.Time          `db:"updated_at"`
	UpdatedBy        nuuid.NUUID        `db:"updated_by" validate:"required"`
	DeletedAt        null.Time          `db:"deleted_at"`
	DeletedBy        nuuid.NUUID        `db:"deleted_by"`
	VariantWarehouse []VariantWarehouse `db:"-" validate:"required,dive,required"`
}

func (pv ProductVariant) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.ToResponseFormat())
}

func (pv ProductVariant) NewFromRequestFormat(format ProductVariantRequestFormat, productID uuid.UUID) (productVariant ProductVariant) {
	variantID, _ := uuid.NewV4()
	productVariant = ProductVariant{
		VariantID:   variantID,
		ProductID:   productID,
		VariantName: format.VariantName,
		Price:       format.Price,
	}
	return
}

func (pv *ProductVariant) ToResponseFormat() ProductVariantResponseFormat {
	resp := ProductVariantResponseFormat{
		VariantID:        pv.VariantID,
		VariantName:      pv.VariantName,
		Price:            pv.Price,
		CreatedAt:        pv.CreatedAt,
		CreatedBy:        pv.CreatedBy,
		UpdatedAt:        pv.UpdatedAt,
		UpdatedBy:        pv.UpdatedBy, // UpdatedBy is already of type nuuid.NUUID (not a pointer)
		DeletedAt:        pv.DeletedAt,
		DeletedBy:        pv.DeletedBy,                              // DeletedBy is already of type nuuid.NUUID (not a pointer)
		VariantWarehouse: make([]VariantWarehouseResponseFormat, 0), // Initialize an empty slice
	}

	// Convert each VariantWarehouse to VariantWarehouseResponseFormat and append to resp.VariantWarehouse.
	for _, variantWarehouse := range pv.VariantWarehouse {
		resp.VariantWarehouse = append(resp.VariantWarehouse, variantWarehouse.ToResponseFormat())
	}

	return resp
}

type ProductVariantRequestFormat struct {
	VariantID        uuid.UUID          `json:"variant_id" validate:"required"`
	VariantName      string             `json:"variant_name" validate:"required"`
	Price            int                `json:"price" validate:"required,min=0"`
	VariantWarehouse []VariantWarehouse `json:"variant_warehouse" validate:"required,dive,required"`
}

type ProductVariantResponseFormat struct {
	VariantID        uuid.UUID          `json:"variant_id"`
	VariantName      string             `json:"variant_name"`
	Price            int                `json:"price"`
	CreatedAt        time.Time          `json:"created_at"`
	CreatedBy        uuid.UUID          `json:"created_by"`
	UpdatedAt        time.Time          `json:"updated_at"`
	UpdatedBy        nuuid.NUUID        `json:"updated_by"`
	DeletedAt        null.Time          `json:"deleted_at"`
	DeletedBy        nuuid.NUUID        `json:"deleted_by"`
	VariantWarehouse []VariantWarehouse `json:"variant_warehouse"`
}

// Variant Warehouse Stock
type VariantWarehouse struct {
	VariantID   uuid.UUID   `db:"variant_id" validate:"required"`
	WarehouseID uuid.UUID   `db:"warehouse_id" validate:"required"`
	Stock       int         `db:"stock" validate:"required"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt   time.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
	Warehouse   []Warehouse `db:"-" validate:"required,dive,required"`
}

func (vw VariantWarehouse) MarshalJSON() ([]byte, error) {
	return json.Marshal(vw.ToResponseFormat())
}

func (vw VariantWarehouse) NewFromRequestFormat(format VariantWarehouseRequestFormat, variantID uuid.UUID, warehouseID uuid.UUID) (variantWarehouse VariantWarehouse) {
	variantWarehouse = VariantWarehouse{
		VariantID:   variantID,
		WarehouseID: warehouseID,
		Stock:       format.Stock,
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
		UpdatedBy:   &vw.UpdatedBy,
		DeletedAt:   vw.DeletedAt,
		DeletedBy:   &vw.DeletedBy,
	}
}

type VariantWarehouseRequestFormat struct {
	VariantID   uuid.UUID `json:"variant_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Stock       int       `json:"stock"`
}

type VariantWarehouseResponseFormat struct {
	VariantID   uuid.UUID    `json:"variant_id"`
	WarehouseID uuid.UUID    `json:"warehouse_id"`
	Stock       int          `json:"stock"`
	CreatedAt   time.Time    `json:"created_at"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	UpdatedAt   time.Time    `json:"updated_at"`
	UpdatedBy   *nuuid.NUUID `json:"updated_by"`
	DeletedAt   null.Time    `json:"deleted_at"`
	DeletedBy   *nuuid.NUUID `json:"deleted_by"`
}

// ProductImage
type ProductImage struct {
	ImageID   uuid.UUID   `db:"image_id" validate:"required"`
	ProductID uuid.UUID   `db:"product_id" validate:"required"`
	ImageURL  string      `db:"image_url"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

func (pi ProductImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(pi.ToResponseFormat())
}

func (pi ProductImage) NewFromRequestFormat(format ProductImageRequestFormat, productID uuid.UUID) (productImage ProductImage) {
	imageID, _ := uuid.NewV4()
	productImage = ProductImage{
		ImageID:  imageID,
		ImageURL: format.ImageURL,
	}

	return
}

func (pi ProductImage) ToResponseFormat() ProductImageResponseFormat {
	return ProductImageResponseFormat{
		ImageID:   pi.ImageID,
		ImageURL:  pi.ImageURL,
		CreatedAt: pi.CreatedAt,
		CreatedBy: pi.CreatedBy,
		UpdatedAt: pi.UpdatedAt,
		UpdatedBy: &pi.UpdatedBy,
		DeletedAt: pi.DeletedAt,
		DeletedBy: &pi.DeletedBy,
	}
}

type ProductImageRequestFormat struct {
	ImageURL string `json:"image_url"`
}

type ProductImageResponseFormat struct {
	ImageID   uuid.UUID    `json:"image_id"`
	ImageURL  string       `json:"image_url"`
	CreatedAt time.Time    `json:"created_at"`
	CreatedBy uuid.UUID    `json:"created_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	UpdatedBy *nuuid.NUUID `json:"updated_by"`
	DeletedAt null.Time    `json:"deleted_at"`
	DeletedBy *nuuid.NUUID `json:"deleted_by"`
}
