package repository

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/db"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
)

type OrderRepo interface {
	Insert(tx db.TX, req model.Order) (id int64, err error)
	GetByID(id int64) (res *model.Order, err error)
	Updata(tx db.TX, req model.Order) (count int64, err error)
	GetHistory(req data.OrderGetHistoryReq) (res []model.OrderGetInfoRes, total int64, err error)
	Delete(tx db.TX, id int64) (count int64, err error)
}

type orderRepo struct {
	pg db.PostgresqlDb
}

func NewOrderRepo(pg db.PostgresqlDb) OrderRepo {
	return &orderRepo{
		pg: pg,
	}
}

func (repo *orderRepo) Insert(tx db.TX, req model.Order) (id int64, err error) {
	params := make([]interface{}, 7)
	params[0] = req.UserID
	params[1] = req.Total
	params[2] = req.Status
	params[3] = req.CreatedBy
	params[4] = req.CreatedDate
	params[5] = req.ModifiedBy
	params[6] = req.ModifiedDate.NullTime

	insertTable := `INSERT INTO orders`

	col := `(
		user_id,
		total,
		status,
		created_by,
		created_date,
		modified_by,
		modified_date
	)`

	values := `VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)`

	returning := `RETURNING id`

	query := repo.pg.Rebind(fmt.Sprintf("%s %s %s %s;", insertTable, col, values, returning))

	err = tx.QueryRow(query, params...).Scan(&id)
	if db.IsSQLReallyError(err) {
		log.Println("SQL insert failed.", err)
		return 0, err
	}

	return id, nil
}

func (repo *orderRepo) GetByID(id int64) (*model.Order, error) {
	sl := `SELECT *`

	from := `FROM orders`

	condition := `WHERE id = ?`

	query := repo.pg.Rebind(fmt.Sprintf("%s %s %s;", sl, from, condition))

	var res model.Order
	err := repo.pg.Get(&res, query, id)
	if db.IsSQLReallyError(err) {
		log.Println("SQL select failed.", err)
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &res, nil

}

func (repo *orderRepo) Updata(tx db.TX, req model.Order) (count int64, err error) {
	params := make([]interface{}, 6)
	params[0] = req.Total
	params[1] = req.Status
	params[2] = req.Reason
	params[3] = req.ModifiedBy
	params[4] = req.ModifiedDate.NullTime
	params[5] = req.ID

	tableUpdate := `UPDATE orders`

	update := `SET
		total = $1,
		status = $2,
		reason = $3,
		modified_by = $4,
		modified_date = $5
	`

	where := `WHERE id = $6`

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s %s;`, tableUpdate, update, where))

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

func (repo *orderRepo) GetHistory(req data.OrderGetHistoryReq) (res []model.OrderGetInfoRes, total int64, err error) {
	params := []interface{}{
		req.UserID,
	}

	sl := `SELECT 
		o.id ,o.user_id ,
		o.status ,o.total ,
		o.reason ,o.modified_date ,
		od.id as order_detail_id,od.product_id ,
		od.quantity ,od.price ,od.sub_total,
		p."name" as product_name ,p.description 
	
	`

	from := `FROM orders o`

	join := `
		LEFT JOIN order_details od ON o.id = od.order_id 
		LEFT JOIN products p ON od.product_id = p.id 
	`

	condition := `WHERE o.user_id = ?`

	conditionStatus := ""
	if req.Status != "" {
		conditionStatus = `AND (o.status = ?)`
		params = append(params, req.Status)
	}

	condition = fmt.Sprintf("%s %s", condition, conditionStatus)

	queryCount := repo.pg.Rebind(fmt.Sprintf(`SELECT COUNT(*) %s %s`, from, condition))

	err = repo.pg.Get(&total, queryCount, params...)
	if db.IsSQLReallyError(err) {
		return nil, 0, err
	}
	if total == 0 {
		return res, 0, nil
	}

	// order
	order := ""
	if req.SortBy.Field != "" {
		if req.SortBy.FieldType == reflect.String {
			order = fmt.Sprintf(`ORDER BY %s COLLATE "th-TH-x-icu" %s, od.id asc`, req.SortBy.Field, req.SortBy.Mode)
		} else {
			order = fmt.Sprintf(`ORDER BY %s %s, od.id asc`, req.SortBy.Field, req.SortBy.Mode)
		}
	}
	// limit
	limit := ""
	if req.PerPage > 0 {
		limit = `LIMIT ?`
		params = append(params, req.PerPage)
	}
	// offset
	offset := ""
	if req.Page > 0 {
		offset = `OFFSET ?`
		params = append(params, (req.Page-1)*req.PerPage)
	}

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s %s %s %s %s %s`,
		sl, from, join,
		condition,
		order, limit, offset,
	))

	err = repo.pg.Select(&res, query, params...)
	if db.IsSQLReallyError(err) {
		return nil, 0, err
	}

	return res, total, nil
}

func (repo *orderRepo) Delete(tx db.TX, id int64) (count int64, err error) {
	params := make([]interface{}, 1)
	params[0] = id

	tableDelete := `DELETE FROM orders`
	condition := `WHERE id = $1`
	query := fmt.Sprintf("%s %s;", tableDelete, condition)

	result, err := tx.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
