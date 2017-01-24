package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/routers"
)

func TestPing(t *testing.T) {
	testRouter := routers.SetupRouter()

	req, err := http.NewRequest("GET", "/v1/ping", nil)
	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
