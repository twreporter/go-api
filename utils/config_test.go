package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenFileError(t *testing.T) {
	// given a non-existed file name
	err := LoadConfig("mocks/confi.json")
	assert.NotNil(t, err)
	assert.Equal(t, err.Message, "utils.config.load_conifg.open_file: ")
}

func TestDecodeFileContentError(t *testing.T) {
	// given a mal form json file
	absFilepath, _ := filepath.Abs("mocks/mal-form-config.json")
	err := LoadConfig(absFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, err.Message, "utils.config.load_config.decode_json: ")
}

func TestLoadConfigSuccess(t *testing.T) {
	absFilepath, _ := filepath.Abs("mocks/config.json")
	err := LoadConfig(absFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, Cfg.EmailSettings)
}
