package mongo

const (
	// Define mongo query operator
	OpAnd          = "$and"
	OpConcatArrays = "$concatArrays"
	OpEq           = "$eq"
	OpExpr         = "$expr"
	OpIn           = "$in"
	OpLet          = "$let"
	OpOr           = "$or"
	OpReduce       = "$reduce"

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
	MetaIn           = "in"
	MetaInitialValue = "initialValue"
	MetaInput        = "input"
	MetaLocalField   = "localField"
	MetaLet          = "let"
	MetaPipeline     = "pipeline"
	MetaVars         = "vars"

	MetaPath                       = "path"
	MetaPreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)
