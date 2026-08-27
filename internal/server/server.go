package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"status": "ok",
		})
	})

	farmerRepo := farmer.NewRepository(db)

	farmerService := farmer.NewService(
		farmerRepo,
		cfg.JWTSecret,
		24*time.Hour,
	)

	farmerHandler := farmer.NewHandler(farmerService)

	router.POST(
		"/api/v1/farmers/register",
		farmerHandler.Register,
	)

	return router
}