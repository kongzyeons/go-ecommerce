package types

import (
	"database/sql"
	"fmt"
)

// SQLNullString sql null string
type SQLNullString struct {
	sql.NullString
}

func NewNullString(value string) SQLNullString {
	val := SQLNullString{}
	if value == "" {
		val.SetNull()
		return val
	}
	val.SetString(value)
	return val
}

// SetString set string to SQLNullString
func (s *SQLNullString) SetString(str string) {
	s.String = str
	s.Valid = true
}

// IsNull SQLNullString isnull
func (s *SQLNullString) IsNull() bool {
	return !s.Valid
}

// SetNull set null to SQLNullString
func (s *SQLNullString) SetNull() {
	s.String = ""
	s.Valid = false
}

// Val get value
func (s SQLNullString) Val() string {
	if s.IsNull() {
		return ""
	}
	return s.String
}

// MarshalJSON Jsonconvert Helper
func (s *SQLNullString) MarshalJSON() ([]byte, error) {

	if s.IsNull() {
		return []byte("\"\""), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", s.Val())), nil
}
