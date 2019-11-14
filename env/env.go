package env

import (
	"os"
)

type ProjectConfig struct {
	SourceDir string
}

func GetConfig() *ProjectConfig {
	return &ProjectConfig{
		SourceDir: getEnv("SOURCE_DIR", "src"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
