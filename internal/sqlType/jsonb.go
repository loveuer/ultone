package sqlType

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
)

type JSONB struct {
	Val   pgtype.JSONB
	Valid bool
}

func NewJSONB(v interface{}) JSONB {
	j := new(JSONB)
	j.Val = pgtype.JSONB{}
	if err := j.Val.Set(v); err == nil {
		j.Valid = true
		return *j
	}

	return *j
}

func (j *JSONB) Set(value interface{}) error {
	if err := j.Val.Set(value); err != nil {
		j.Valid = false
		return err
	}

	j.Valid = true

	return nil
}

func (j *JSONB) Bind(model interface{}) error {
	return j.Val.AssignTo(model)
}

func (j *JSONB) Scan(value interface{}) error {
	j.Val = pgtype.JSONB{}
	if value == nil {
		j.Valid = false
		return nil
	}

	j.Valid = true

	return j.Val.Scan(value)
}

func (j JSONB) Value() (driver.Value, error) {
	if j.Valid {
		return j.Val.Value()
	}

	return nil, nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return j.Val.MarshalJSON()
	}

	return json.Marshal(nil)
}

func (j *JSONB) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		j.Valid = false
		return j.Val.UnmarshalJSON(b)
	}

	return j.Val.UnmarshalJSON(b)
}
