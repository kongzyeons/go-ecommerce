package repository

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/db"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type ProductRepo interface {
	GetList(req data.ProductGetListReq) (res []model.Product, total int64, err error)
	GetPriceByIDs(ids ...int64) (res []model.GetPriceByIDsRes, err error)
}

type productRepo struct {
	pg db.PostgresqlDb
}

func NewProductRepo(pg db.PostgresqlDb) ProductRepo {
	return &productRepo{
		pg: pg,
	}
}

func (repo *productRepo) GetList(req data.ProductGetListReq) (res []model.Product, total int64, err error) {
	params := []interface{}{}

	sl := `SELECT *`

	from := `FROM products`

	condition := `WHERE TRUE`

	conditionSearch := ""
	if req.SearchText != "" {
		words := strings.Fields(req.SearchText)
		condtionWord := make([]string, len(words))
		for i, word := range words {
			condtionWord[i] = `(name ILIKE ?)`
			params = append(params, `%\`+word+"%")
		}
		conditionSearch = `AND ( ` + strings.Join(condtionWord, " AND ") + ` )`
	}

	condition = fmt.Sprintf("%s %s", condition, conditionSearch)

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
			order = fmt.Sprintf(`ORDER BY %s COLLATE "th-TH-x-icu" %s`, req.SortBy.Field, req.SortBy.Mode)
		} else {
			order = fmt.Sprintf(`ORDER BY %s %s`, req.SortBy.Field, req.SortBy.Mode)
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

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s %s %s %s %s`,
		sl, from,
		condition,
		order, limit, offset,
	))

	err = repo.pg.Select(&res, query, params...)
	if db.IsSQLReallyError(err) {
		return nil, 0, err
	}

	return res, total, nil
}

func (repo *productRepo) GetPriceByIDs(ids ...int64) (res []model.GetPriceByIDsRes, err error) {
	withReq := `WITH req AS (
	SELECT *
	FROM UNNEST($1::bigint[])
		AS t(id)
	)`

	sl := `SELECT 
		products.id,
		products.price
	`

	from := `FROM products`

	join := `RIGHT JOIN req ON req.id = products.id`

	query := repo.pg.Rebind(fmt.Sprintf(`%s %s %s %s;`, withReq, sl, from, join))

	err = repo.pg.Select(&res, query, pq.Array(ids))
	if db.IsSQLReallyError(err) {
		return res, err
	}

	return res, nil
}
