package sqlType

import (
	"database/sql"
	"encoding/json"
)

// type NullString struct {
// 	sql.NullString
// }

type NullString struct{ sql.NullString }

func NewNullString(val string) NullString {
	if val == "" {
		return NullString{}
	}

	return NullString{sql.NullString{Valid: true, String: val}}
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &ns.String); err != nil {
		ns.Valid = true
		return err
	}

	return nil
}
