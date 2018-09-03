package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"

	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

// Register recieves http POST requests and create the registration record in the storage
func (mc *MembershipController) Register(c *gin.Context) {
	var err error
	var appErr models.AppError
	var postBody models.RegistrationJSON
	var registration models.Registration
	var service string
	var activeToken string

	postBody, err = mc.parseRegPOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	service = c.Param("service")

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)
	if err != nil {
		c.JSON(err.(models.AppError).StatusCode, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	postBody.ActivateToken = activeToken

	registration, err = mc.Storage.CreateRegistration(service, postBody)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message,
			"error": err.Error()})
		return
	}

	registration.ActivateToken = ""

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "record": registration})

	// TODO Send activation email
}

// Deregister recieves http DELETE request and delete the registration recode in the storage
func (mc *MembershipController) Deregister(c *gin.Context) {
	var err error
	var userEmail, service string
	var appErr models.AppError

	userEmail = c.Param("userEmail")
	service = c.Param("service")

	err = mc.Storage.DeleteRegistration(userEmail, service)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// GetRegisteredUser recieves http GET request and get the registration record from the storage
func (mc *MembershipController) GetRegisteredUser(c *gin.Context) {
	var appErr models.AppError

	email := c.Param("userEmail")
	service := c.Param("service")

	reg, err := mc.Storage.GetRegistration(email, service)

	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	// not to return actviate token
	reg.ActivateToken = ""

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": reg})
	return
}

// GetRegisteredUsers recieves the http GET request which may contain limit, offset, order_by and active_code params, and get regitration records from the storage
func (mc *MembershipController) GetRegisteredUsers(c *gin.Context) {
	var appErr models.AppError
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

	count, err = mc.Storage.GetRegistrationsAmountByService(service, activeCode)
	registrations, err = mc.Storage.GetRegistrationsByService(service, offset, limit, orderBy, activeCode)

	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message,
			"error": err.Error()})
		return
	}

	for i, registration := range registrations {
		registration.ActivateToken = ""
		registrations[i] = registration
	}

	c.JSON(200, gin.H{"status": "ok", "count": count, "records": registrations})
}

// ActivateRegistration recieves http GET request, and make the registration record active
func (mc *MembershipController) ActivateRegistration(c *gin.Context) {
	email := c.Param("userEmail")
	service := c.Param("service")
	token := c.Query("activeToken")

	u := url.URL{
		Host:   viper.GetString("consumersettings.host"),
		Scheme: viper.GetString("consumersettings.protocol"),
		Path:   constants.Activate,
	}

	reg, err := mc.Storage.GetRegistration(email, service)

	if err == nil {
		if reg.ActivateToken == token {
			_, err = mc.Storage.UpdateRegistration(service, models.RegistrationJSON{Email: email, Active: true, ActivateToken: token})
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

func (mc *MembershipController) parseRegPOSTBody(c *gin.Context) (models.RegistrationJSON, error) {
	var err error
	var json models.RegistrationJSON

	contentType := c.ContentType()

	if contentType == "application/json" {
		err = c.Bind(&json)
		if err != nil {
			return models.RegistrationJSON{}, models.NewAppError("getPOSTBody", "Bad request", err.Error(), http.StatusBadRequest)
		}
		return json, nil
	}

	return models.RegistrationJSON{}, models.NewAppError("getPOSTBody", "Bad request", "POST body is not JSON formatted", http.StatusBadRequest)
}
