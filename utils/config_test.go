package utils

import (
	"strings"
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func TestOpenFileError(t *testing.T) {
	// given a non-existed file name
	err := LoadConfig("mocks/confi.json")
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "utils.config.load_conifg.open_file: "))
}

func TestDecodeFileContentError(t *testing.T) {
	// given a mal form json file
	absFilepath, _ := filepath.Abs("mocks/mal-form-config.json")
	err := LoadConfig(absFilepath)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "utils.config.load_config.decode_json: "))
}

func TestLoadConfigSuccess(t *testing.T) {
	absFilepath, _ := filepath.Abs("mocks/config.json")
	err := LoadConfig(absFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, Cfg.EmailSettings)
}
