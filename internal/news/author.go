package news

import (
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/pkg/errors"
)

func GetAuthorWithIndex(index AlgoliaSearcher, q *Query) ([]Author, int, error) {
	var authors []Author
	res, err := index.Search(q.Filter.Name, opt.Offset(q.Offset), opt.Length(q.Limit))
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
	return authors, res.NbHits, nil
}
