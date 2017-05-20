package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
)

// ServiceController defines the routes and methods to handle the requests
type ServiceController struct {
	Storage storage.MembershipStorage
}

// SetRoute set the route path and corresponding handlers
func (sc ServiceController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {

	group.POST("/services", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), sc.Create)
	group.DELETE("/services/:name", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), sc.Delete)
	group.PUT("/services/:name", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), sc.Update)
	group.GET("/services/:name", middlewares.CheckJWT(), sc.Read)

	return group
}

// Create recieves http POST requests and create the service record in the storage
func (sc ServiceController) Create(c *gin.Context) {
	var err error
	var appErr models.AppError
	var postBody models.ServiceJSON
	var service models.Service

	postBody, err = sc.parsePOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	service, err = sc.Storage.CreateService(postBody)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": service})
}

// Delete recieves http DELETE request and delete the service record in the storage
func (sc ServiceController) Delete(c *gin.Context) {

	name := c.Param("name")

	err := sc.Storage.DeleteService(name)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// Read recieves http GET request and get the service record in the storage
func (sc ServiceController) Read(c *gin.Context) {

	name := c.Param("name")

	svc, err := sc.Storage.GetService(name)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": svc})
	return
}

// Update recieves http PUT request and update/create the service record in the storage
func (sc ServiceController) Update(c *gin.Context) {
	var err error
	var appErr models.AppError
	var postBody models.ServiceJSON
	var service models.Service

	name := c.Param("name")

	postBody, err = sc.parsePOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	service, err = sc.Storage.UpdateService(name, postBody)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": service})
}

func (sc ServiceController) parsePOSTBody(c *gin.Context) (models.ServiceJSON, error) {
	var err error
	var json models.ServiceJSON

	contentType := c.ContentType()

	if contentType == "application/json" {
		err = c.Bind(&json)
		if err != nil {
			return models.ServiceJSON{}, models.NewAppError("getPOSTBody", "Bad request", err.Error(), http.StatusBadRequest)
		}
		return json, nil
	}

	return models.ServiceJSON{}, models.NewAppError("getPOSTBody", "Bad request", "POST body is not JSON formatted", http.StatusBadRequest)
}
