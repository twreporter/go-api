package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
)

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestMain(m *testing.M) {
	var err error
	// Create DB connections
	if DB, err = OpenGormConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	}

	DB.SetJoinTableHandler(&models.User{}, constants.TableBookmarks, &models.UsersBookmarks{})

	// Create Mongo DB connections
	if MgoDB, err = OpenMgoConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to mongo database, but got err=%+v", err))
	}
	// Set up database, including drop existing tables and create tables
	RunMigration()

	// Set up default records in tables
	SetDefaultRecords()

	SetupGinServer()

	retCode := m.Run()
	os.Exit(retCode)
}
