package model

import (
	"app-ecommerce/pkg/types"
	"time"
)

type Order struct {
	ID           int64               `db:"id"`
	UserID       int64               `db:"user_id"`
	Total        float64             `db:"total"`
	Status       string              `db:"status"`
	Reason       types.SQLNullString `db:"reason"`
	CreatedBy    string              `db:"created_by"`
	CreatedDate  time.Time           `db:"created_date"`
	ModifiedBy   string              `db:"modified_by"`
	ModifiedDate types.SQLNullTime   `db:"modified_date"`
}

type OrderGetInfoRes struct {
	ID           int64               `db:"id"`
	UserID       int64               `db:"user_id"`
	Total        float64             `db:"total"`
	Status       string              `db:"status"`
	Reason       types.SQLNullString `db:"reason"`
	ModifiedDate types.SQLNullTime   `db:"modified_date"`

	// join order detail
	OrderDetailID types.SQLNullInt64   `db:"order_detail_id"`
	ProductID     types.SQLNullInt64   `db:"product_id"`
	Quantity      types.SQLNullInt64   `db:"quantity"`
	Price         types.SQLNullFloat64 `db:"price"`
	SubTotal      types.SQLNullFloat64 `db:"sub_total"`

	// join product
	ProductName types.SQLNullString `db:"product_name"`
	Description types.SQLNullString `db:"description"`
}
