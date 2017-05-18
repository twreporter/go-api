package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"

	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

// RegistrationController defines the routes and methods to handle the requests
type RegistrationController struct {
	Storage storage.MembershipStorage
}

// SetRoute set the route path and corresponding handlers
func (rc RegistrationController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {

	// TODO add middleware to check the request from twreporter.org domain
	group.POST("/registrations/:service", rc.Register)

	// TODO add middleware to check the email to delete is the email of the user sending the request
	group.DELETE("/registrations/:service/:userEmail", rc.Deregister)

	// TODO add middleware to check the request from twreporter.org domain
	group.GET("/registrations/:service/:userEmail", rc.GetRegisteredUser)
	group.GET("/registrations/:service", rc.GetRegisteredUsers)
	//

	group.GET("/activation/:service/:userEmail", rc.Activate)

	return group
}

// Register recieves http POST requests and create the registration record in the storage
func (rc RegistrationController) Register(c *gin.Context) {
	var err error
	var postBody models.RegistrationJSON
	var registration models.Registration
	var service string
	var activeToken string

	postBody, err = rc.parsePOSTBody(c)
	if err != nil {
		log.Error("controllers.registration.register.error_to_parse_post_body: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	service = c.Param("service")

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)
	if err != nil {
		log.Error("controllers.registration.register.error_to_generate_active_token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	postBody.Service = service
	postBody.ActivateToken = activeToken

	registration, err = rc.Storage.CreateRegistration(postBody)
	if err != nil {
		log.Error("controllers.registration.register.error_to_create_registration: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	registration.ActivateToken = ""

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": registration})

	// TODO Send activation email
}

// Deregister recieves http DELETE request and delete the registration recode in the storage
func (rc RegistrationController) Deregister(c *gin.Context) {
	var err error
	var userEmail, service string

	userEmail = c.Param("userEmail")
	service = c.Param("service")

	err = rc.Storage.DeleteRegistration(userEmail, service)
	if err != nil {
		log.Error("controllers.registration.register.error_to_delete_registration: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// GetRegisteredUser recieves http GET request and get the registration record from the storage
func (rc RegistrationController) GetRegisteredUser(c *gin.Context) {
	email := c.Param("userEmail")
	service := c.Param("service")

	reg, err := rc.Storage.GetRegistration(email, service)

	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		log.Error("controllers.registration.get_registered_user.error_to_get: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"status": "Resource not found", "error": err.Error()})
		return
	} else if err != nil {
		log.Error("controllers.registration.get_registered_user.error_to_get: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	// not to return actviate token
	reg.ActivateToken = ""

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": reg})
	return
}

// GetRegisteredUsers recieves the http GET request which may contain limit, offset, order_by and active_code params, and get regitration records from the storage
func (rc RegistrationController) GetRegisteredUsers(c *gin.Context) {
	var registrations []models.Registration
	var limit, offset, activeCode int
	var count uint
	var err error

	service := c.Param("service")

	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = constants.DefaultLimit
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = constants.DefaultOffset
	}
	activeCode, err = strconv.Atoi(c.Query("active_code"))
	if err != nil {
		activeCode = constants.DefaultActiveCode
	}
	orderBy := c.Query("order_by")

	if orderBy == "" {
		orderBy = constants.DefaultOrderBy
	}

	count, err = rc.Storage.GetRegistrationsAmountByService(service, activeCode)
	registrations, err = rc.Storage.GetRegistrationsByService(service, offset, limit, orderBy, activeCode)

	if err != nil {
		log.Error("controllers.registration.get_registered_users.error_to_get_registrations_from_storage: ", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	for i, registration := range registrations {
		registration.ActivateToken = ""
		registrations[i] = registration
	}

	c.JSON(200, gin.H{"status": "ok", "count": count, "records": registrations})
}

// Activate recieves http GET request, and make the registration record active
func (rc RegistrationController) Activate(c *gin.Context) {
	email := c.Param("userEmail")
	service := c.Param("service")
	token := c.Query("activeToken")

	u := url.URL{
		Host:   utils.Cfg.ConsumerSettings.Host,
		Scheme: utils.Cfg.ConsumerSettings.Protocal,
		Path:   constants.Activate,
	}

	reg, err := rc.Storage.GetRegistration(email, service)

	if err == nil {
		if reg.ActivateToken == token {
			_, err = rc.Storage.UpdateRegistration(models.RegistrationJSON{Email: email, Service: service, Active: true, ActivateToken: token})
			if err == nil {
				c.Redirect(http.StatusTemporaryRedirect, u.String())
				return
			}
		} else {
			q := u.Query()
			q.Set("error", "Activate token is not correct")
			q.Set("error_code", strconv.Itoa(http.StatusForbidden))
			u.RawQuery = q.Encode()
			c.Redirect(http.StatusTemporaryRedirect, u.String())
			return
		}
	}

	q := u.Query()
	q.Set("error", err.Error())
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusTemporaryRedirect, u.String())
	return

}

func (rc RegistrationController) parsePOSTBody(c *gin.Context) (models.RegistrationJSON, error) {
	var err error
	var json models.RegistrationJSON

	contentType := c.ContentType()

	if contentType == "application/json" {
		err = c.Bind(&json)
		if err != nil {
			return models.RegistrationJSON{}, models.NewAppError("getPOSTBody", "controllers.registration_controller.parse_post_body", err.Error(), http.StatusBadRequest)
		}
		return json, nil
	}

	return models.RegistrationJSON{}, models.NewAppError("getPOSTBody", "controllers.registration_controller.parse_post_body", "POST body is not JSON formatted", http.StatusBadRequest)
}
