package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/database"
	"github.com/raaj2493/KrishiSetu/internal/server"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	cfg := config.Load()

	_, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := server.New()
	log.Printf("KrishiSetu server starting on :%s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}