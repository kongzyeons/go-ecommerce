package repository

import (
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/db"
	"fmt"
	"log"
)

type UserRepo interface {
	Insert(tx db.TX, req model.User) (id int64, err error)
	GetCountUnique(req model.User) (count int, err error)
}

type userRepo struct {
	pg db.PostgresqlDb
}

func NewUserRepo(pg db.PostgresqlDb) UserRepo {
	return &userRepo{
		pg: pg,
	}
}

func (repo *userRepo) Insert(tx db.TX, req model.User) (id int64, err error) {
	params := make([]interface{}, 8)
	params[0] = req.Name
	params[1] = req.Email
	params[2] = req.Password
	params[3] = req.RoleID
	params[4] = req.CreatedBy
	params[5] = req.CreatedDate
	params[6] = req.ModifiedBy
	params[7] = req.ModifiedDate.NullTime

	insertTable := `INSERT INTO users`

	col := `(
		name,
		email,
		password,
		role_id,
		created_by,
		created_date,
		modified_by,
		modified_date
	)`
	values := `VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
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

func (repo *userRepo) GetCountUnique(req model.User) (count int, err error) {
	var params []interface{}

	sl := `SELECT COUNT(*)`

	from := `FROM users`

	condition := `WHERE TRUE`

	condtionWithoutID := ""
	if req.ID > 0 {
		condtionWithoutID = `AND (id != ?)`
		params = append(params, req.ID)
	}

	condtionName := ""
	if req.Name != "" {
		condtionName = `AND (name = ?)`
		params = append(params, req.Name)
	}

	conditionEmail := ""
	if req.Email != "" {
		conditionEmail = `AND (email = ?)`
		params = append(params, req.Email)
	}

	condition = fmt.Sprintf("%s %s %s %s", condition, condtionWithoutID, condtionName, conditionEmail)

	query := repo.pg.Rebind(fmt.Sprintf("%s %s %s;", sl, from, condition))

	err = repo.pg.Get(&count, query, params...)
	if db.IsSQLReallyError(err) {
		log.Println("SQL select failed.", err)
		return 0, err
	}

	return count, nil
}
