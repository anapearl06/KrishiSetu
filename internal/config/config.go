package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv             string
	AppPort            string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours time.Duration
}

func Load() Config {
	expirationHours := 24

	if value := os.Getenv("JWT_EXPIRATION_HOURS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			expirationHours = parsed
		}
	}

	return Config{
		AppEnv:             os.Getenv("APP_ENV"),
		AppPort:            os.Getenv("APP_PORT"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: time.Duration(expirationHours) * time.Hour,
	}
}