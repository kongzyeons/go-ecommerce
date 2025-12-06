package repository

import (
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/db"
	"fmt"
	"strings"
)

type OrderDetailRepo interface {
	InsertMany(tx db.TX, orderID int64, req []model.OrderDetail) (count int64, err error)
	DeleteMany(tx db.TX, orderID int64) (count int64, err error)
}

type orderDetailRepo struct {
	pg db.PostgresqlDb
}

func NewOrderDetailRepo(pg db.PostgresqlDb) OrderDetailRepo {
	return &orderDetailRepo{
		pg: pg,
	}
}

func (repo *orderDetailRepo) InsertMany(tx db.TX, orderID int64, req []model.OrderDetail) (count int64, err error) {
	if len(req) == 0 {
		return 0, nil
	}

	insertTable := `INSERT INTO order_details`

	col := `(
		order_id,
		product_id,
		quantity,
		price,
		sub_total
	)`

	params := []interface{}{}
	valueStrings := []string{}
	for _, v := range req {
		base := len(params)
		ph := make([]string, 0, 5)
		for i := 1; i <= 5; i++ {
			ph = append(ph, fmt.Sprintf("$%d", base+i))
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(ph, ",")))
		params = append(params,
			orderID,
			v.ProductID,
			v.Quantity,
			v.Price,
			v.SubTotal,
		)
	}

	values := fmt.Sprintf(`VALUES %s`, strings.Join(valueStrings, ","))

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s %s;`, insertTable, col, values))

	res, err := tx.Exec(query, params...)
	if err != nil {
		return 0, err
	}
	count, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *orderDetailRepo) DeleteMany(tx db.TX, orderID int64) (count int64, err error) {
	tableDelete := `DELETE FROM order_details`

	condition := `WHERE order_id = $1`

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s;`, tableDelete, condition))

	result, err := tx.Exec(query, orderID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
