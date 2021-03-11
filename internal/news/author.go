package news

import (
	"context"
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/pkg/errors"
)

func GetAuthorWithIndex(ctx context.Context, index AlgoliaSearcher, q *Query) ([]Author, int64, error) {
	var authors []Author
	res, err := index.Search(q.Filter.Name, opt.Offset(q.Offset), opt.Length(q.Limit), ctx)
	if err != nil {
		// fallback
		return nil, -1, errors.WithStack(err)
	}
	rawRecords, err := json.Marshal(res.Hits)
	if err != nil {
		// fallback
		return nil, -1, errors.WithStack(err)
	}
	err = json.Unmarshal(rawRecords, &authors)
	if err != nil {
		// fallback
		return nil, -1, errors.WithStack(err)
	}
	return authors, int64(res.NbHits), nil
}
