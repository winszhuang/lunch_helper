package util

import "database/sql"

func CheckNullString(text string) sql.NullString {
	result := sql.NullString{}
	if text == "" {
		result.Valid = false
		result.String = ""
	} else {
		result.Valid = true
		result.String = text
	}
	return result
}

func CheckNullInt32(value int) sql.NullInt32 {
	result := sql.NullInt32{}
	if value == 0 {
		result.Valid = false
		result.Int32 = 0
	} else {
		result.Valid = true
		result.Int32 = int32(value)
	}
	return result
}
