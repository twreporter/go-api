package news

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/internal/query"
)

func TestParseSinglePostQuery(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want *Query
	}{
		{
			name: "Given default parameter",
			url:  "http://example.com/slug",
			want: &Query{
				Filter: Filter{Slug: "slug"},
			},
		},
		{
			name: "Given the full parameter",
			url:  "http://example.com/slug?full=true",
			want: &Query{
				Filter: Filter{Slug: "slug"},
				Full:   true,
			},
		},
		{
			name: "Given query parameters, ignore unsupported",
			url:  "http://example.com/slug?full=true&unsupported=value",
			want: &Query{
				Filter: Filter{Slug: "slug"},
				Full:   true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := helperCreateContext(t, tc.url, "/:slug")
			got := ParseSinglePostQuery(c)
			if !reflect.DeepEqual(*got, *tc.want) {
				t.Errorf("expected query %v, got %v", tc.want, got)
			}
		})
	}
}

func TestParsePostListQuery(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want *Query
	}{
		{
			name: "Given default parameter",
			url:  "http://example.com/posts",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given the category_id parameter",
			url:  "http://example.com/posts?category_id=cid1&category_id=cid2",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Filter:     Filter{Categories: []string{"cid1", "cid2"}},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given the tag_id parameter",
			url:  "http://example.com/posts?tag_id=tid1&tag_id=tid2",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Filter:     Filter{Tags: []string{"tid1", "tid2"}},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given the id parameter",
			url:  "http://example.com/posts?id=id1&id=id2",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Filter:     Filter{Tags: []string{"id1", "id2"}},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given the sort by published_date parameter",
			url:  "http://example.com/posts?sort=published_date",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(true)}},
			},
		},
		{
			name: "Given the sort by updated_at parameter",
			url:  "http://example.com/posts?sort=-updated_at",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Sort:       SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given the pagination parameters",
			url:  "http://example.com/posts?offset=5&limit=6",
			want: &Query{
				Pagination: query.Pagination{Offset: 5, Limit: 6},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
		{
			name: "Given query parameters, ignore unsupported one",
			url:  "http://example.com/posts?unsupported=value",
			want: &Query{
				Pagination: query.Pagination{Offset: 0, Limit: 10},
				Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := helperCreateContext(t, tc.url, "/posts")
			got := ParsePostListQuery(c)
			if !reflect.DeepEqual(*got, *tc.want) {
				t.Errorf("expected query %v, got %v", tc.want, got)
			}
		})
	}
}

func TestParseSingleTopicQuery(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want *Query
	}{
		{
			name: "Given default parameter",
			url:  "http://example.com/topics/slug",
			want: &Query{
				Filter: Filter{Slug: "slug"},
			},
		},
		{
			name: "Given the full parameter",
			url:  "http://example.com/topics/slug?full=true",
			want: &Query{
				Filter: Filter{Slug: "slug"},
				Full:   true,
			},
		},
		{
			name: "Given query parameters, ignore unsupported",
			url:  "http://example.com/topics/slug?full=true&unsupported=value",
			want: &Query{
				Filter: Filter{Slug: "slug"},
				Full:   true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := helperCreateContext(t, tc.url, "/topics/:slug")
			got := ParseSingleTopicQuery(c)
			if !reflect.DeepEqual(*got, *tc.want) {
				t.Errorf("expected query %v, got %v", tc.want, got)
			}
		})
	}
}

func helperCreateContext(t *testing.T, url, route string) *gin.Context {
	t.Helper()
	c, router := gin.CreateTestContext(httptest.NewRecorder())
	// Placeholder route handler for matching the targeting url
	router.GET(route, func(*gin.Context) {})
	// Set the request to force the context rebuild as our desired one
	c.Request = httptest.NewRequest(http.MethodGet, url, nil)
	router.HandleContext(c)
	return c
}
