package db

import (
	"context"
	"database/sql"
)

type TX interface {
	Rollback() error
	Commit() error
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) SqlRow
	QueryRowContext(ctx context.Context, query string, args ...any) SqlRow
}

type SqlRow interface {
	Err() error
	Scan(dest ...any) error
}

type transactionX struct {
	tx *sql.Tx
}

func (t *transactionX) Rollback() error {
	return t.tx.Rollback()
}
func (t *transactionX) Commit() error {
	return t.tx.Commit()
}
func (t *transactionX) Exec(query string, args ...any) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *transactionX) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *transactionX) QueryRow(query string, args ...any) SqlRow {
	return t.tx.QueryRow(query, args...)
}

func (t *transactionX) QueryRowContext(ctx context.Context, query string, args ...any) SqlRow {
	return t.tx.QueryRowContext(ctx, query, args...)
}
