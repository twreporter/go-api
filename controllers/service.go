package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

// ServiceController defines the routes and methods to handle the requests
type ServiceController struct {
	Storage storage.MembershipStorage
}

// SetRoute set the route path and corresponding handlers
func (sc ServiceController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {

	// TODO add middleware to check permission
	group.POST("/services", sc.Create)

	// TODO add middleware to check permission
	group.DELETE("/services/:id", sc.Delete)

	// TODO add middleware to check permission
	group.PUT("/services/:id", sc.Update)

	// TODO add middleware to check permission
	group.GET("/services/:id", sc.Read)

	return group
}

// Create recieves http POST requests and create the service record in the storage
func (sc ServiceController) Create(c *gin.Context) {
	var err error
	var postBody models.ServiceJSON
	var service models.Service

	postBody, err = sc.parsePOSTBody(c)
	if err != nil {
		log.Error("controllers.service.create.error_to_parse_post_body: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	service, err = sc.Storage.CreateService(postBody)
	if err != nil {
		log.Error("controllers.service.register.error_to_create_service: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": service})

	// TODO Send activation email
}

// Delete recieves http DELETE request and delete the service record in the storage
func (sc ServiceController) Delete(c *gin.Context) {
	var err error
	var id string

	id = c.Param("id")

	err = sc.Storage.DeleteService(id)
	if err != nil {
		log.Error("controllers.service.delete.error_to_delete_service: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// Read recieves http GET request and get the service record in the storage
func (sc ServiceController) Read(c *gin.Context) {

	id := c.Param("id")

	svc, err := sc.Storage.GetService(id)

	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		log.Error("controllers.service.get_service.error_to_get: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "Resource not found", "error": err.Error()})
		return
	} else if err != nil {
		log.Error("controllers.service.get_service.error_to_get: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": svc})
	return
}

// Update recieves http PUT request and update/create the service record in the storage
func (sc ServiceController) Update(c *gin.Context) {
	var idInt int
	var err error
	var postBody models.ServiceJSON
	var service models.Service

	id := c.Param("id")

	postBody, err = sc.parsePOSTBody(c)
	if err != nil {
		log.Error("controllers.service.create.error_to_parse_post_body: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	service, err = sc.Storage.UpdateService(id, postBody)
	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		idInt, err = strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
			return
		}
		postBody.ID = uint(idInt)
		service, err = sc.Storage.CreateService(postBody)
		c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": service})
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
			return models.ServiceJSON{}, models.NewAppError("getPOSTBody", "controllers.service_controller.parse_post_body", err.Error(), http.StatusBadRequest)
		}
		return json, nil
	}

	return models.ServiceJSON{}, models.NewAppError("getPOSTBody", "controllers.service_controller.parse_post_body", "POST body is not JSON formatted", http.StatusBadRequest)
}
