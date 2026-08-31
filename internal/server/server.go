package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// CORS middleware for Vercel frontend
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"status": "ok",
		})
	})

	// =========================
	// Farmer
	// =========================

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

		farmerRoutes.PUT(
			"/me",
			middleware.JWTAuth(cfg.JWTSecret),
			farmerHandler.UpdateProfile,
		)
	}

	// =========================
	// Buyer
	// =========================

	buyerRepo := buyer.NewRepository(db)

	buyerService := buyer.NewService(
		buyerRepo,
		cfg.JWTSecret,
		cfg.JWTExpirationHours,
	)

	buyerHandler := buyer.NewHandler(buyerService)

	buyerRoutes := router.Group("/api/v1/buyers")
	{
		buyerRoutes.POST(
			"/register",
			buyerHandler.Register,
		)

		buyerRoutes.POST(
			"/login",
			buyerHandler.Login,
		)

		buyerRoutes.GET(
			"/me",
			middleware.JWTAuth(cfg.JWTSecret),
			buyerHandler.Me,
		)

		buyerRoutes.PUT(
			"/me",
			middleware.JWTAuth(cfg.JWTSecret),
			buyerHandler.UpdateProfile,
		)
	}

	// =========================
	// Listings (frontend-aligned)
	// =========================

	listingRepo := listing.NewRepository(db)
	listingService := listing.NewService(listingRepo)
	listingHandler := listing.NewHandler(listingService)

	listingRoutes := router.Group("/api/v1/listings")
	{
		// Public: browse marketplace
		listingRoutes.GET(
			"",
			listingHandler.Browse,
		)

		// Protected: farmer listing management (JWT)
		listingRoutes.POST(
			"",
			middleware.JWTAuth(cfg.JWTSecret),
			listingHandler.Create,
		)

		listingRoutes.GET(
			"/my",
			middleware.JWTAuth(cfg.JWTSecret),
			listingHandler.GetMyListings,
		)

		listingRoutes.PUT(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			listingHandler.Update,
		)

		listingRoutes.DELETE(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			listingHandler.Delete,
		)
	}

	return router
}
