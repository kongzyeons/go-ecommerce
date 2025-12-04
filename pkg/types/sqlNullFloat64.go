package types

import (
	"database/sql"
	"strconv"

	"github.com/shopspring/decimal"
)

// SQLNullFloat64 sql null float64
type SQLNullFloat64 struct {
	sql.NullFloat64
}

func NewNullFloat64(value float64) SQLNullFloat64 {
	val := SQLNullFloat64{}
	val.SetFloat64(value)
	return val
}

// IsNull check null
func (s SQLNullFloat64) IsNull() bool {
	return !s.Valid
}

// Val get int64
func (s SQLNullFloat64) Val(defaultVal float64) float64 {
	if s.IsNull() {
		return defaultVal
	}
	return s.Float64
}

// SetFloat64 set float64
func (s *SQLNullFloat64) SetFloat64(val float64) {
	s.Float64 = val
	s.Valid = true
}

// SetNull set null
func (s *SQLNullFloat64) SetNull() {
	s.Float64 = 0.0
	s.Valid = false
}

// SetDecimal set decimal
func (s *SQLNullFloat64) SetDecimal(dec decimal.Decimal) {
	f, _ := dec.Float64()
	s.SetFloat64(f)
}

// Decimal to decimal
func (s *SQLNullFloat64) Decimal() decimal.Decimal {
	return decimal.NewFromFloat(s.Val(0.0))
}

func (s *SQLNullFloat64) String() string {

	val := strconv.FormatFloat(s.Val(0.0), 'f', 2, 64)

	if val == "0.00" {
		return ""
	}

	return val

}
