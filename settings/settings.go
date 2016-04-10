package settings

import (
	"os"
)

// Type settings type
type Type struct {
	ElasticSearchHost string
	ElasticSearchPort string
	TestServer        string
}

func getVarFromEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// LoadSettings load settings
func LoadSettings() Type {
	Settings := Type{}
	Settings.ElasticSearchHost = getVarFromEnv("ENV_ELASTICSEARCH_HOST", "localhost")
	Settings.ElasticSearchPort = getVarFromEnv("ENV_ELASTICSEARCH_PORT", "9200")
	Settings.TestServer = getVarFromEnv("ENV_TEST_SERVER", "http://localhost:8001")
	return Settings
}
