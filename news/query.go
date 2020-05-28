package news

import "gopkg.in/guregu/null.v4"

type (
	OrderBy struct {
		IsAsc null.Bool
	}

	Pagination struct {
		Offset uint
		Limit  uint
	}

	PostFilter struct {
		Slug string
	}

	PostQuery struct {
		Pagination
		Filter PostFilter
		Sort   PostSort
	}

	PostSort struct {
		UpdatedAt OrderBy
	}

	TopicFilter struct {
		Slug string
	}

	TopicQuery struct {
		Pagination
		Filter TopicFilter
		Sort   TopicSort
	}

	TopicSort struct {
		UpdatedAt OrderBy
	}
)
