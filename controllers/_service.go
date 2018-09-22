package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"twreporter.org/go-api/models"
)

// Create recieves http POST requests and create the service record in the storage
func (mc *MembershipController) Create(c *gin.Context) {
	var err error
	var appErr models.AppError
	var postBody models.ServiceJSON
	var service models.Service

	postBody, err = mc.parseServicePOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	service, err = mc.Storage.CreateService(postBody)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": service})
}

// Delete recieves http DELETE request and delete the service record in the storage
func (mc *MembershipController) Delete(c *gin.Context) {

	name := c.Param("name")

	err := mc.Storage.DeleteService(name)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// Read recieves http GET request and get the service record in the storage
func (mc *MembershipController) Read(c *gin.Context) {

	name := c.Param("name")

	svc, err := mc.Storage.GetService(name)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": svc})
	return
}

// Update recieves http PUT request and update/create the service record in the storage
func (mc *MembershipController) Update(c *gin.Context) {
	var err error
	var appErr models.AppError
	var postBody models.ServiceJSON
	var service models.Service

	name := c.Param("name")

	postBody, err = mc.parseServicePOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	service, err = mc.Storage.UpdateService(name, postBody)
	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": service})
}

func (mc *MembershipController) parseServicePOSTBody(c *gin.Context) (models.ServiceJSON, error) {
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
