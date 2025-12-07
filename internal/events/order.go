package events

import (
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/types"
	"time"
)

type OrderEvent struct {
	Role       string  `json:"role"`
	SendID     int64   `json:"sendID"`
	OrderID    int64   `json:"orderID"`
	Status     string  `json:"status"`
	UserID     int64   `json:"userID"`
	Total      float64 `json:"total"`
	Reason     string  `json:"reason"`
	ModifiedBy string  `json:"modifiedBy"`
}

func (obj *OrderEvent) ToOrderUpdateDB() model.Order {
	return model.Order{
		ID:           obj.OrderID,
		Status:       obj.Status,
		UserID:       obj.UserID,
		Total:        obj.Total,
		Reason:       types.NewNullString(obj.Reason),
		ModifiedBy:   obj.ModifiedBy,
		ModifiedDate: types.NewNullTime(time.Now().UTC()),
	}
}
