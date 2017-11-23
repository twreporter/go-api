package utils

import (
	"database/sql"

	"twreporter.org/go-api/configs/constants"

	log "github.com/Sirupsen/logrus"
)

// ToNullString invalidates a sql.NullString if empty, validates if not empty
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// GetNullString returns a invalid NullString
func GetNullString() sql.NullString {
	return sql.NullString{String: "", Valid: false}
}

// GetGender format the gender string
func GetGender(s string) sql.NullString {
	var ngender sql.NullString
	switch s {
	case "":
		ngender = GetNullString()
	case "male":
		ngender = ToNullString(constants.GenderMale)
	case "female":
		ngender = ToNullString(constants.GenderFemale)
	default:
		// Other gender
		ngender = ToNullString(constants.GenderOthers)
	}
	return ngender
}

func Check(f func() error) {
	if err := f(); err != nil {
		log.Error("Received error:", err.Error())
	}
}
