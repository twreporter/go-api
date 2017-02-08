package configs

import "github.com/jinzhu/configor"

// Config ...
type Config struct {
	DB struct {
		Name     string `required:"true"`
		User     string `default:"root"`
		Password string `required:"true"`
		Address  string `default:"127.0.0.1"`
		Port     string `default:"3306"`
	}
}

// GetConfig ...
func GetConfig() Config {
	// load the config for the API
	var config = Config{}
	configor.Load(&config, "config.yml")

	return config
}
