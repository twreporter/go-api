package mongo

const (
	// Define mongo query operator
	OpIn = "$in"

	OrderAsc  = 1
	OrderDesc = -1

	// Define mongo pipeline stage
	StageLimit = "$limit"
	StageMatch = "$match"
	StageSkip  = "$skip"
	StageSort  = "$sort"
)
