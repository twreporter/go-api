// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"twreporter.org/go-api/models"
)

// Cfg it is used to store the data of config file
var Cfg = &models.Config{}

// CfgFileName it is filename of config file
var CfgFileName string

func init() {
	Cfg.SetDefaults()
}

// LoadConfig it will load config file
func LoadConfig(fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		appError := models.NewAppError("LoadConfig", "internal server error: fail to load config", err.Error(), http.StatusInternalServerError)
		return appError
	}

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		appError := models.NewAppError("LoadConfig", "internal server error: fail to load config", err.Error(), http.StatusInternalServerError)
		return appError
	}

	config.SetDefaults()

	Cfg = &config

	return nil
}
