package types

import (
	"database/sql"
	"fmt"
	"strconv"
)

// SQLNullInt64 sql type for in64
type SQLNullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(value int64) SQLNullInt64 {
	val := SQLNullInt64{}
	if value == 0 {
		val.SetNull()
		return val
	}
	val.SetInt64(value)
	return val
}

// IsNull check null
func (s SQLNullInt64) IsNull() bool {
	return !s.Valid
}

// Val get int64
func (s SQLNullInt64) Val(defaultVal int64) int64 {
	if s.IsNull() {
		return defaultVal
	}
	return s.Int64
}

// SetInt64 set int64
func (s *SQLNullInt64) SetInt64(val int64) {
	s.Int64 = val
	s.Valid = true
}

func (s *SQLNullInt64) String() string {

	val := strconv.FormatInt(s.Int64, 10)

	if val == "0" {
		return ""
	}

	return val
}

// SetNull set null
func (s *SQLNullInt64) SetNull() {
	s.Int64 = 0
	s.Valid = false
}

// MarshalJSON Jsonconvert Helper
func (s *SQLNullInt64) MarshalJSON() ([]byte, error) {

	if s.IsNull() {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("%d", s.Val(0))), nil
}

// GetInt64OrNull - get int or nil
func (s SQLNullInt64) GetIntOrNull() *int64 {

	if s.IsNull() {
		return nil
	}

	return &s.Int64

}

// GetInt64 - get int
func (s SQLNullInt64) GetInt() int64 {
	if s.IsNull() {
		return 0
	}
	return s.Int64
}
