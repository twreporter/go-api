package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"twreporter.org/go-api/storage"
)

func toResponse(err error) (int, gin.H, error) {
	cause := errors.Cause(err)

	// For legacy storage errors, the NotFound and Conflict error type is responsed with status "error" if no explicit handlers in controller layer.
	// Try to migrate these errors to status "fail".
	// TODO: adjust client side errors
	switch {
	case storage.IsNotFound(err):
		return http.StatusNotFound, gin.H{"status": "error", "message": fmt.Sprintf("record not found. %s", cause.Error())}, nil
	case storage.IsConflict(err):
		return http.StatusConflict, gin.H{"status": "error", "message": fmt.Sprintf("record is already existed. %s", cause.Error())}, nil
	default:
		// omit itentionally
	}

	return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("internal server error. %s", cause.Error())}, nil
}

func toPostResponse(err error) (int, gin.H, error) {
	cause := errors.Cause(err)

	// For legacy post storage errors, payload contained `status`:"[Custom Error Message]" and `error`:"[Detail message]"
	// if no explicit handlers in controller layer specified.
	// Try to migrate these errors to `status`:"fail" or "error" and replace `error`field with corresponding `data` or `message` field
	// adhered to jsend payload.
	// TODO: adjust error payloads
	switch {
	case storage.IsNotFound(err):
		return http.StatusNotFound, gin.H{"status": fmt.Sprintf("record not found. %s", cause.Error()), "error": cause.Error()}, nil
	case storage.IsConflict(err):
		return http.StatusConflict, gin.H{"status": fmt.Sprintf("record is already existed. %s", cause.Error()), "error": cause.Error()}, nil
	default:
		// omit itentionally
	}

	return http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("internal server error. %s", cause.Error()), "error": cause.Error()}, nil
}
