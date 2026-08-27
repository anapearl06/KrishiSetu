package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv        string
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() Config {
	expirationHours, err := strconv.Atoi(
		os.Getenv("JWT_EXPIRATION_HOURS"),
	)

	if err != nil || expirationHours <= 0 {
		expirationHours = 24
	}

	return Config{
		AppEnv:        os.Getenv("APP_ENV"),
		AppPort:       os.Getenv("APP_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: time.Duration(expirationHours) * time.Hour,
	}
}