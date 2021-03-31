package news

import "github.com/algolia/algoliasearch-client-go/v3/algolia/search"

type AlgoliaSearcher interface {
	Search(query string, opts ...interface{}) (res search.QueryRes, err error)
}
