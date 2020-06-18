package news

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
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
