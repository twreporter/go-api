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
	StageFilter    = "$filter"
	StageLimit     = "$limit"
	StageLookup    = "$lookup"
	StageMatch     = "$match"
	StageSkip      = "$skip"
	StageSort      = "$sort"
	StageUnwind    = "$unwind"
	StageProject   = "$project"

	// Define Meta fields for nested stages (e.g., lookup)
	MetaAs           = "as"
	MetaCond         = "cond"
	MetaForeignField = "foreignField"
	MetaFrom         = "from"
	MetaInput        = "input"
	MetaLocalField   = "localField"
	MetaLet          = "let"
	MetaPipeline     = "pipeline"

	MetaPath                       = "path"
	MetaPreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)
