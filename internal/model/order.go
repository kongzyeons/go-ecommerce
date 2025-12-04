package model

import (
	"app-ecommerce/pkg/types"
	"time"
)

type Order struct {
	ID           int64             `db:"id"`
	UserID       int64             `db:"user_id"`
	Total        float64           `db:"total"`
	Status       string            `db:"status"`
	CreatedBy    string            `db:"created_by"`
	CreatedDate  time.Time         `db:"created_date"`
	ModifiedBy   string            `db:"modified_by"`
	ModifiedDate types.SQLNullTime `db:"modified_date"`
}
