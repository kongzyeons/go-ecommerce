package model

import "app-ecommerce/pkg/types"

type Product struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

type GetPriceByIDsRes struct {
	ID    types.SQLNullInt64   `json:"id"`
	Price types.SQLNullFloat64 `json:"price"`
}
