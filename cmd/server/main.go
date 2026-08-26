package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/raaj2493/KrishiSetu/internal/config"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	cfg := config.Load()

	fmt.Println("Environment:", cfg.AppEnv)
	fmt.Println("Port:", cfg.AppPort)
	fmt.Println("Database:", cfg.DBName)
}