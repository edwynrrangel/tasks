package config

import "os"

type Configuration struct {
	Port               string
	CorsAllowedOrigins string
}

func LoadConfig() *Configuration {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsAllowedOrigins == "" {
		corsAllowedOrigins = "*"
	}

	return &Configuration{
		Port:               port,
		CorsAllowedOrigins: corsAllowedOrigins,
	}
}

var Config *Configuration
