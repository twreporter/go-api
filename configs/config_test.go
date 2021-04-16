package configs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twreporter/go-api/configs"
)

func TestLoadConf(t *testing.T) {
	// If environment variable is provided,
	// it should overwrite the default config
	t.Run("Environment variables overwrite default", func(t *testing.T) {

		const (
			testEnvFirstLevel          = "test"
			testEnvNestedConfig        = "test_protocol"
			testEnvNestedUnderscoreKey = "test_issuer"
			testEnvNestedArray         = "http://testhost1 http://testhost2"
		)

		// First-level config
		os.Setenv("GOAPI_ENVIRONMENT", testEnvFirstLevel)
		// Nested config
		os.Setenv("GOAPI_APP_PROTOCOL", testEnvNestedConfig)
		// Nested config with underscore key
		os.Setenv("GOAPI_APP_JWT_ISSUER", testEnvNestedUnderscoreKey)
		// Nested config with array of values
		os.Setenv("GOAPI_CORS_ALLOW_ORIGINS", testEnvNestedArray)

		testConf, _ := configs.LoadConf("")

		//Validate output
		assert.Equal(t, testConf.Environment, testEnvFirstLevel)
		assert.Equal(t, testConf.App.Protocol, testEnvNestedConfig)
		assert.Equal(t, testConf.App.JwtIssuer, testEnvNestedUnderscoreKey)
		assert.Equal(t, testConf.Cors.AllowOrigins, []string{
			"http://testhost1",
			"http://testhost2",
		})
	})
}
