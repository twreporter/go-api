package news

import (
	"context"
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/pkg/errors"
)

// GetRankedAuthorIDs returns ranked author ID result by index search
func GetRankedAuthorIDs(ctx context.Context, index AlgoliaSearcher, q *Query) ([]string, int64, error) {
	type authorIndex struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var indexes []authorIndex
	res, err := index.Search(q.Filter.Name, opt.Offset(q.Offset), opt.Length(q.Limit), ctx)
	if err != nil {
		// fallback
		return nil, 0, errors.WithStack(err)
	}
	rawRecords, err := json.Marshal(res.Hits)
	if err != nil {
		// fallback
		return nil, 0, errors.WithStack(err)
	}
	err = json.Unmarshal(rawRecords, &indexes)
	if err != nil {
		// fallback
		return nil, 0, errors.WithStack(err)
	}
	var authorIDs []string
	for _, index := range indexes {
		authorIDs = append(authorIDs, index.ID)
	}
	return authorIDs, int64(res.NbHits), nil
}
