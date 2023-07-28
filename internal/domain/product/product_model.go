package product

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Product struct {
	ProductID   uuid.UUID   `db:"product_id" validate:"required"`
	BrandID     uuid.UUID   `db:"brand_id" validate:"required"`
	UserID      uuid.UUID   `db:"user_id" validate:"required"`
	ProductName string      `db:"product_name"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt   time.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
	Variants    []Variant   `db:"-" validate:"required,dive,required"`
}

type Variant struct {
	VariantID   uuid.UUID   `db:"variant_id" validate:"required"`
	ProductID   uuid.UUID   `db:"product_id" validate:"required"`
	VariantName string      `db:"variant_name"`
	Price       int         `db:"price" validate:"required"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt   time.Time   `db:"updated_at"`
	UpdatedBy   nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt   null.Time   `db:"deleted_at"`
	DeletedBy   nuuid.NUUID `db:"deleted_by"`
	Images      []Image     `db:"-" validate:"required,dive,required"`
	Warehouse   []Warehouse `db:"-" validate:"required,dive,required"`
}

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

type Image struct {
	ImageID   uuid.UUID   `db:"image_id" validate:"required"`
	VariantID uuid.UUID   `db:"variant_id" validate:"required"`
	ImageURL  string      `db:"image_url"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
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

func (p *Product) IsDeleted() (deleted bool) {
	return p.DeletedAt.Valid && p.DeletedBy.Valid
}

func (p *Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}

func (p Product) NewFromRequestFormat(req ProductRequestFormat, userID uuid.UUID) (newProduct Product, err error) {

	productID, _ := uuid.NewV4()
	newProduct = Product{
		ProductID: productID,
		// brandID can be implemented by passing variable
		BrandID:     brandID,
		UserID:      userID,
		ProductName: req.ProductName,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
		// UpdatedAt using trigger
		// UpdatedBy using trigger
		// Implement signature later
	}

	variants := make([]Variant, 0)
	for _, requestVariant := range req.Variants {
		variant := Variant{}
		variant = variant.NewFromRequestFormat(requestVariant, productID)
		variants = append(variants, variant)
	}

	err = newProduct.Validate()

	return
}

func (p Product) ToResponseFormat() ProductResponseFormat {
	resp := ProductResponseFormat{}
}

func (p *Product) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(p)
}

type ProductRequestFormat struct {
	Name     string `json:"name" validate:"required"`
	ShopName string `json:"shopName" validate:"required"`
	Price    int    `json:"price" validate:"required, min=0"`
}

type ProductResponseFormat struct {
	Name      string    `json:"name"`
	ShopName  string    `json:"shopName"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt,omitempty"`
}
