// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"encoding/json"
	"os"
	"twreporter.org/go-api/models"
)

// Cfg it is used to store the data of config file
var Cfg *models.Config = &models.Config{}

// CfgFileName it is filename of config file
var CfgFileName string

// LoadConfig will try to search around for the corresponding config file.
// It will search ../config/fileName and last it will look at fileName
func LoadConfig(fileName string) {

	file, err := os.Open("configs/" + fileName)
	if err != nil {
		panic("utils.config.load_config.open.panic: " + err.Error())
	}

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic("utils.config.load_config.decoding.panic: " + err.Error())
	}

	config.SetDefaults()

	Cfg = &config
}
