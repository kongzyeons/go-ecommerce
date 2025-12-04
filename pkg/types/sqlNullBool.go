package types

import (
	"database/sql"
	"strconv"
)

// SQLNullBool SQL null Boolean
type SQLNullBool struct {
	sql.NullBool
}

func NewNullBool(value bool) SQLNullBool {
	val := SQLNullBool{}
	val.SetBool(value)
	return val
}

// IsNull check null
func (s SQLNullBool) IsNull() bool {
	return !s.Valid
}

// Val get bool value
func (s SQLNullBool) Val() bool {
	if s.IsNull() {
		return false
	}
	return s.Bool
}

// SetInt64 set int64
func (s *SQLNullBool) SetBool(val bool) {
	s.Bool = val
	s.Valid = true
}

func (s *SQLNullBool) String() string {
	val := strconv.FormatBool(s.Bool)
	// if val == "0" {
	// 	return ""
	// }
	return val
}

// SetNull set null
func (s *SQLNullBool) SetNull() {
	s.Bool = false
	s.Valid = false
}
