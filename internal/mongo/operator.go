package mongo

const (
	// Define mongo query operator
	OpAnd          = "$and"
	OpConcatArrays = "$concatArrays"
	OpEq           = "$eq"
	OpExpr         = "$expr"
	OpIn           = "$in"
	OpLet          = "$let"
	OpGte          = "$gte"
	OpOr           = "$or"
	OpReduce       = "$reduce"
	OpExists       = "$exists"
	OpNe           = "$ne"
	OpCount        = "$count"
	OpNot          = "$not"
	OpSize         = "$size"

	OrderAsc  = 1
	OrderDesc = -1

	ElemMatch = "$elemMatch"

	// Define mongo pipeline stage
	StageAddFields   = "$addFields"
	StageGroup       = "$group"
	StageFilter      = "$filter"
	StageLimit       = "$limit"
	StageLookup      = "$lookup"
	StageMatch       = "$match"
	StageSkip        = "$skip"
	StageSort        = "$sort"
	StageUnwind      = "$unwind"
	StageReplaceRoot = "$replaceRoot"
	StageProject     = "$project"
	StageFacet       = "$facet"

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
