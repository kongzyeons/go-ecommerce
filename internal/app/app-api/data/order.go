package data

type OrderDetail struct {
	ProductID int64 `json:"productID" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required"`
}

type OrderCreateReq struct {
	OrderDetails []OrderDetail `json:"orderDetails" validate:"required"`
	UserID       int64         `json:"-" validate:"required"`
	CreateBy     string        `json:"-" validate:"required"`
}
