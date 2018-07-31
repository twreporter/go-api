package utils

import (
	"net/http"
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
)

func TestOpenFileError(t *testing.T) {
	// given a non-existed file name
	err := LoadConfig("mocks/confi.json")
	assert.NotNil(t, err)

	appErr := err.(*models.AppError)
	assert.Equal(t, appErr.Where, "LoadConfig")
	assert.Equal(t, appErr.StatusCode, http.StatusInternalServerError)
}

func TestDecodeFileContentError(t *testing.T) {
	// given a mal form json file
	absFilepath, _ := filepath.Abs("mocks/mal-form-config.json")
	err := LoadConfig(absFilepath)
	assert.NotNil(t, err)

	appErr := err.(*models.AppError)
	assert.Equal(t, appErr.Where, "LoadConfig")
	assert.Equal(t, appErr.StatusCode, http.StatusInternalServerError)
}

func TestLoadConfigSuccess(t *testing.T) {
	absFilepath, _ := filepath.Abs("mocks/config.json")
	err := LoadConfig(absFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, Cfg.EmailSettings)
}
