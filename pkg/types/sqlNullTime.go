package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

const mysqlDatetimeFormat = "2006-01-02 15:04:05"
const defaultTimeZone = ""

// SQLNullTime - null Time
type SQLNullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time) SQLNullTime {
	if t.IsZero() {
		return SQLNullTime{
			NullTime: sql.NullTime{
				Time:  t,
				Valid: false,
			},
		}
	}
	return SQLNullTime{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: true,
		},
	}
}

// Scan implements the Scanner interface.
func (t *SQLNullTime) Scan(value interface{}) error {

	if val, ok := value.(time.Time); ok {

		//logger.Debug.Printf("%+v", val.Location())

		t.Time = val
		t.Valid = true

	} else if val, ok := value.([]byte); ok {

		// logger.Debug.Printf("%v\n", string(val))
		// logger.Debug.Printf("%v\n", time.Local.String())

		ti, err := time.ParseInLocation(mysqlDatetimeFormat, string(val), time.UTC)
		if err != nil {
			t.Time = time.Time{}
			t.Valid = false
			return nil
		}
		t.Time = ti
		t.Valid = true

		// logger.Debug.Printf("Local location: %v\n", ti.Format("2006-01-02 15:04:05 -0700"))
		// logger.Debug.Printf("UTC location: %v\n", ti.UTC())

	}
	return nil

}

// Value implements the driver Valuer interface.
func (t SQLNullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// IsNull check null
func (t SQLNullTime) IsNull() bool {
	return !t.Valid
}

// Val return value
func (t SQLNullTime) Val() time.Time {
	return t.Time
}

// Local return value
func (t SQLNullTime) Local() time.Time {
	return t.Time.Local()
}

// SetTime set time
func (t *SQLNullTime) SetTime(set time.Time) {
	t.Scan(set)
}

// DateString - create date string format with spit text ex. 01/01/2018
func (t SQLNullTime) DateString(spit string) string {

	if t.Valid == false {
		return ""
	}

	time := t.Val()
	dateString := fmt.Sprintf("%02d%s%02d%s%04d", time.Day(), spit, time.Month(), spit, time.Year())

	return dateString

}

// String - Set Time to String Date
func (t SQLNullTime) String() string {

	if t.Valid == false {
		return ""
	}

	dateString := t.Val().String()

	return dateString
}

// GetTimeOrNull - get time or nil
func (t SQLNullTime) GetTimeOrNull() *time.Time {

	if t.IsNull() {
		return nil
	}

	return &t.Time

}

// MarshalJSON Jsonconvert Helper
func (s *SQLNullTime) MarshalJSON() ([]byte, error) {

	if s.IsNull() {
		return []byte("\"\""), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", s.Val().Format("2006-01-02 15:04:05"))), nil
}
