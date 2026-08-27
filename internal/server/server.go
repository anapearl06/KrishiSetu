package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
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
		cfg.JWTExpirationHours,
	)

	farmerHandler := farmer.NewHandler(farmerService)

	farmerRoutes := router.Group("/api/v1/farmers")
	{
		farmerRoutes.POST(
			"/register",
			farmerHandler.Register,
		)

		farmerRoutes.POST(
			"/login",
			farmerHandler.Login,
		)

		farmerRoutes.GET(
			"/me",
			middleware.JWTAuth(cfg.JWTSecret),
			farmerHandler.Me,
		)
	}

	return router
}