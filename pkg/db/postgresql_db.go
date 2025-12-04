package db

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

type PostgresqlDb interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecTx(ctx context.Context, fn func(TX) error) error
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) SqlRow
	Exec(query string, args ...any) (sql.Result, error)
	Rebind(query string) string
}

type postgresqlDb struct {
}

var postgresDBInstance PostgresqlDb
var postgresDBOnce sync.Once

func NewPostgresqlDb() PostgresqlDb {
	postgresDBOnce.Do(func() {
		postgresDBInstance = postgresqlDb{}
	})
	return postgresDBInstance
}

func (pg postgresqlDb) Get(dest interface{}, query string, args ...interface{}) error {
	return dbx.Get(dest, query, args...)
}

func (pg postgresqlDb) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dbx.GetContext(ctx, dest, query, args...)
}

func (pg postgresqlDb) Select(dest interface{}, query string, args ...interface{}) error {
	return dbx.Select(dest, query, args...)
}

func (pg postgresqlDb) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dbx.SelectContext(ctx, dest, query, args...)
}

func (pg postgresqlDb) ExecTx(ctx context.Context, fn func(TX) error) error {
	tx, err := pg.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		log.Println("SQL Commit failed.", err)
		return err
	}

	return tx.Commit()
}

func (pg postgresqlDb) BeginTx(ctx context.Context, opts *sql.TxOptions) (TX, error) {
	tx, err := dbx.BeginTx(ctx, opts)
	return &transactionX{tx: tx}, err
}

func (pg postgresqlDb) Query(query string, args ...any) (*sql.Rows, error) {
	return dbx.Query(query, args...)
}

func (pg postgresqlDb) QueryRow(query string, args ...any) SqlRow {
	return dbx.QueryRow(query, args...)
}
func (pg postgresqlDb) Exec(query string, args ...any) (sql.Result, error) {
	return dbx.Exec(query, args...)
}

func (pg postgresqlDb) Rebind(query string) string {
	return dbx.Rebind(query)
}
