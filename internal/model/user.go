package model

import (
	"app-ecommerce/pkg/types"
	"time"
)

type User struct {
	ID           int64             `db:"id"`
	Name         string            `db:"name"`
	Email        string            `db:"email"`
	Password     string            `db:"password"`
	RoleID       int64             `db:"role_id"`
	CreatedBy    string            `db:"created_by"`
	CreatedDate  time.Time         `db:"created_date"`
	ModifiedBy   string            `db:"modified_by"`
	ModifiedDate types.SQLNullTime `db:"modified_date"`
	IsDeleted    bool              `db:"is_deleted"`
}

type UserRole struct {
	ID       int64               `db:"id"`
	Name     string              `db:"name"`
	Email    string              `db:"email"`
	Password string              `db:"password"`
	RoleID   int64               `db:"role_id"`
	RoleName types.SQLNullString `db:"role_name"`
}
