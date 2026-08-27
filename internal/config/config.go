package config

import "os"

type Config struct {
	AppEnv             string
	AppPort            string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours string
}

func Load() Config {
	return Config{
		AppEnv:             os.Getenv("APP_ENV"),
		AppPort:            os.Getenv("APP_PORT"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: os.Getenv("JWT_EXPIRATION_HOURS"),
	}
}