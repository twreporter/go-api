package utils

import (
	"gopkg.in/guregu/null.v3"

	"github.com/twreporter/go-api/configs/constants"
)

// GetGender format the gender string
func GetGender(s string) null.String {
	var gender string
	switch s {
	case "":
		gender = s
	case "male":
		gender = constants.GenderMale
	case "female":
		gender = constants.GenderFemale
	default:
		// Other gender
		gender = constants.GenderOthers
	}
	return null.StringFrom(gender)
}
