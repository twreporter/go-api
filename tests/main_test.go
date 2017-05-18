package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/controllers"
)

var (
	Engine *gin.Engine
	DB     *gorm.DB
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
	if DB, err = OpenTestConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	}

	// Set up database, including drop existing tables and create tables
	RunMigration(DB)

	// Set up default records in tables
	SetDefaultRecords(DB)

	cf := controllers.NewControllerFactory(DB)

	Engine = gin.Default()
	routerGroup := Engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	retCode := m.Run()
	os.Exit(retCode)
}
