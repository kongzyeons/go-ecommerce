package model

type OrderDetail struct {
	ID        int64   `db:"id"`
	OrderID   int64   `db:"order_id"`
	ProductID int64   `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
	SubTotal  float64 `db:"sub_total"`
}
