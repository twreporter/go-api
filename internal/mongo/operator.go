package mongo

const (
	// Define mongo query operator
	OpIn = "$in"

	OrderAsc  = 1
	OrderDesc = -1

	// Define mongo pipeline stage
	StageLimit  = "$limit"
	StageLookup = "$lookup"
	StageMatch  = "$match"
	StageSkip   = "$skip"
	StageSort   = "$sort"
	StageUnwind = "$unwind"

	// Define meta fields for nested stages (e.g., lookup)
	metaAs           = "as"
	metaForeignField = "foreignField"
	metaFrom         = "from"
	metaLocalField   = "localField"

	metaPath                       = "path"
	metaPreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)
