package query

// package query provides common query data type for go-api

import "gopkg.in/guregu/null.v3"

type Pagination struct {
	Offset int
	Limit  int
}

type Order struct {
	IsAsc null.Bool
}
