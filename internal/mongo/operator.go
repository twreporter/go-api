package mongo

const (
	// Define mongo query operator
	OpAnd  = "$and"
	OpEq   = "$eq"
	OpExpr = "$expr"
	OpIn   = "$in"

	OrderAsc  = 1
	OrderDesc = -1

	// Define mongo pipeline stage
	StageAddFields = "$addFields"
	StageLimit     = "$limit"
	StageLookup    = "$lookup"
	StageMatch     = "$match"
	StageSkip      = "$skip"
	StageSort      = "$sort"
	StageUnwind    = "$unwind"
	StageProject   = "$project"

	// Define Meta fields for nested stages (e.g., lookup)
	MetaAs           = "as"
	MetaForeignField = "foreignField"
	MetaFrom         = "from"
	MetaLocalField   = "localField"
	MetaLet          = "let"
	MetaPipeline     = "pipeline"

	MetaPath                       = "path"
	MetaPreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)
