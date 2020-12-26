package config

import "os"

func GetEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// we can use viper here to handle configs
