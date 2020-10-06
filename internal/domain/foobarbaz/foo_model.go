package foobarbaz

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// FooStatus indicates the status of Foo.
type FooStatus string

const (
	// FooStatusNew indicates a new Foo.
	FooStatusNew FooStatus = "new"
	// FooStatusPending indicates a Foo that requires verification.
	FooStatusPending FooStatus = "pending"
	// FooStatusVerified indicates a Foo that passed verification.
	FooStatusVerified FooStatus = "verified"
	// FooStatusPaid indicates a Foo that has been paid.
	FooStatusPaid FooStatus = "paid"
	// FooStatusInTransit indicates a Foo that has left our warehouse.
	FooStatusInTransit FooStatus = "inTransit"
	// FooStatusDelivered indicates a Foo that is in the customer's possession.
	FooStatusDelivered FooStatus = "delivered"
	// FooStatusFailedToDeliver indicates a Foo that isn't delivered and is
	// back in our possession.
	FooStatusFailedToDeliver FooStatus = "failedToDeliver"
)

//// Foo

// Foo is a sample parent entity model.
type Foo struct {
	ID            uuid.UUID   `db:"entity_id" validate:"required"`
	Name          string      `db:"name" validate:"required"`
	TotalQuantity int64       `db:"total_quantity" validate:"required,min=1"`
	TotalPrice    float64     `db:"total_price" validate:"required,min=0"`
	TotalDiscount float64     `db:"total_discount" validate:"required,min=0"`
	ShippingFee   float64     `db:"shipping_fee" validate:"required,min=0"`
	GrandTotal    float64     `db:"grand_total" validate:"required,min=0"`
	Status        FooStatus   `db:"status" validate:"required"`
	Created       time.Time   `db:"created" validate:"required"`
	CreatedBy     uuid.UUID   `db:"created_by" validate:"required"`
	Updated       null.Time   `db:"updated"`
	UpdatedBy     nuuid.NUUID `db:"updated_by"`
	Items         []FooItem   `db:"-" validate:"required,dive,required"`
}

// NewFromRequestFormat creates a new Foo from its request format.
func (f Foo) NewFromRequestFormat(req FooRequestFormat, userID uuid.UUID) (newFoo Foo) {
	fooID, _ := uuid.NewV4()
	newFoo = Foo{
		ID:          fooID,
		Name:        req.Name,
		ShippingFee: req.ShippingFee,
		Status:      req.Status,
		Created:     time.Now(),
		CreatedBy:   userID,
	}

	requestItems := make([]FooItem, 0)
	for _, itemRequest := range req.Items {
		item := FooItem{}
		item = item.NewFromRequestFormat(itemRequest, fooID)
		requestItems = append(requestItems, item)
	}
	newFoo.Items = requestItems

	newFoo.Recalculate()

	return
}

// AttachItems attaches FooItems to this Foo.
func (f *Foo) AttachItems(items []FooItem) Foo {
	for _, item := range items {
		if item.FooID == f.ID {
			f.Items = append(f.Items, item)
		}
	}
	return *f
}

// Recalculate recalculates totals in this Foo.
func (f *Foo) Recalculate() {
	f.TotalQuantity = int64(0)
	f.TotalDiscount = float64(0)
	f.TotalPrice = float64(0)
	recalculatedItems := make([]FooItem, 0)
	for _, item := range f.Items {
		item.Recalculate()
		recalculatedItems = append(recalculatedItems, item)
		f.TotalQuantity += item.Quantity
		f.TotalDiscount += item.Discount
		f.TotalPrice += item.TotalPrice
	}
	f.Items = recalculatedItems
	f.GrandTotal = f.TotalPrice - f.TotalDiscount + f.ShippingFee
}

// Validate validates the entity.
func (f *Foo) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(f)
}

// ToResponseFormat converts this Foo to its response format.
func (f Foo) ToResponseFormat() FooResponseFormat {
	resp := FooResponseFormat{
		ID:            f.ID,
		Name:          f.Name,
		TotalQuantity: f.TotalQuantity,
		TotalPrice:    f.TotalPrice,
		TotalDiscount: f.TotalDiscount,
		ShippingFee:   f.ShippingFee,
		GrandTotal:    f.GrandTotal,
		Status:        f.Status,
		Created:       f.Created,
		CreatedBy:     f.CreatedBy,
		Updated:       f.Updated,
		UpdatedBy:     f.UpdatedBy.Ptr(),
		Items:         make([]FooItemResponseFormat, 0),
	}

	for _, item := range f.Items {
		resp.Items = append(resp.Items, item.ToResponseFormat())
	}

	return resp
}

// MarshalJSON overrides the standard JSON formatting.
func (f Foo) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.ToResponseFormat())
}

// FooRequestFormat represents a Foo's standard formatting for JSON deserializing.
type FooRequestFormat struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	ShippingFee float64                `json:"shippingFee"`
	Status      FooStatus              `json:"status"`
	Items       []FooItemRequestFormat `json:"items"`
}

// FooResponseFormat represents a Foo's standard formatting for JSON serializing.
type FooResponseFormat struct {
	ID            uuid.UUID               `json:"id"`
	Name          string                  `json:"name"`
	TotalQuantity int64                   `json:"totalQuantity"`
	TotalPrice    float64                 `json:"totalPrice"`
	TotalDiscount float64                 `json:"totalDiscount"`
	ShippingFee   float64                 `json:"shippingFee"`
	GrandTotal    float64                 `json:"grandTotal"`
	Status        FooStatus               `json:"status"`
	Created       time.Time               `json:"created"`
	CreatedBy     uuid.UUID               `json:"createdBy"`
	Updated       null.Time               `json:"updated,omitempty"`
	UpdatedBy     *uuid.UUID              `json:"updatedBy,omitempty"`
	Items         []FooItemResponseFormat `json:"items"`
}

//// Foo Item

// FooItem is a sample child entity model.
type FooItem struct {
	ID          uuid.UUID `db:"entity_id" validate:"required"`
	FooID       uuid.UUID `db:"foo_id" validate:"required"`
	SKU         string    `db:"sku" validate:"required"`
	ProductName string    `db:"product_name" validate:"required"`
	Quantity    int64     `db:"quantity" validate:"required,min=1"`
	UnitPrice   float64   `db:"unit_price" validate:"required,min=0"`
	TotalPrice  float64   `db:"total_price" validate:"required,min=0"`
	Discount    float64   `db:"discount" validate:"required,min=0"`
	GrandTotal  float64   `db:"grand_total" validate:"required,min=0"`
}

// NewFromRequestFormat creates a new FooItem from its request format.
func (fi FooItem) NewFromRequestFormat(format FooItemRequestFormat, fooID uuid.UUID) (fooItem FooItem) {
	fooItemID, _ := uuid.NewV4()
	fooItem = FooItem{
		ID:          fooItemID,
		FooID:       fooID,
		SKU:         format.SKU,
		ProductName: format.ProductName,
		Quantity:    format.Quantity,
		UnitPrice:   format.UnitPrice,
		Discount:    format.Discount,
	}
	return
}

// Recalculate recalculates totals in this FooItem.
func (fi *FooItem) Recalculate() {
	fi.TotalPrice = float64(fi.Quantity) * fi.UnitPrice
	fi.GrandTotal = fi.TotalPrice - fi.Discount
}

// ToResponseFormat converts this FooItem to its response format.
func (fi *FooItem) ToResponseFormat() FooItemResponseFormat {
	return FooItemResponseFormat{
		ID:          fi.ID,
		FooID:       fi.FooID,
		SKU:         fi.SKU,
		ProductName: fi.ProductName,
		Quantity:    fi.Quantity,
		UnitPrice:   fi.UnitPrice,
		TotalPrice:  fi.TotalPrice,
		Discount:    fi.Discount,
		GrandTotal:  fi.GrandTotal,
	}
}

// MarshalJSON overrides the standard JSON formatting.
func (fi FooItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(fi.ToResponseFormat())
}

// FooItemRequestFormat represents a FooItem's standard formatting for JSON deserializing.
type FooItemRequestFormat struct {
	ID          uuid.UUID `json:"id"`
	FooID       uuid.UUID `json:"fooId"`
	SKU         string    `json:"sku"`
	ProductName string    `json:"productName"`
	Quantity    int64     `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
	Discount    float64   `json:"discount"`
}

// FooItemResponseFormat represents a FooItem's standard formatting for JSON serializing.
type FooItemResponseFormat struct {
	ID          uuid.UUID `json:"entityId"`
	FooID       uuid.UUID `json:"fooId"`
	SKU         string    `json:"sku"`
	ProductName string    `json:"productName"`
	Quantity    int64     `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
	TotalPrice  float64   `json:"totalPrice"`
	Discount    float64   `json:"discount"`
	GrandTotal  float64   `json:"grandTotal"`
}
