package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/offer"
	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/order"
	"github.com/raaj2493/KrishiSetu/internal/server/response"

	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// =========================
	// CORS
	// =========================

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://krishisetuio.vercel.app",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"status": "ok",
		})
	})

	// API version 1 root group
	api := router.Group("/api/v1")

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

	farmerRoutes := api.Group("/farmers")
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

	buyerRoutes := api.Group("/buyers")
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
	// Listings
	// =========================

	listingRepo := listing.NewRepository(db)

	listingService := listing.NewService(listingRepo)

	listingHandler := listing.NewHandler(listingService)

	listings := api.Group("/listings")

	listings.Use(middleware.JWTAuth(cfg.JWTSecret))

	{
		listings.POST("", listingHandler.CreateListing)

		listings.GET("", listingHandler.ListListings)

		listings.GET("/my", listingHandler.GetMyListings)

		listings.GET("/:id", listingHandler.GetListing)

		listings.PUT("/:id", listingHandler.UpdateListing)

		listings.DELETE("/:id", listingHandler.CancelListing)
	}

	offerRepo := offer.NewRepository(db)
	offerService := offer.NewService(offerRepo, db)
	offerHandler := offer.NewHandler(offerService)

	offers := api.Group("/offers")
	offers.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		offers.POST("", offerHandler.CreateOffer)
		offers.GET("/my", offerHandler.GetMyOffers)
		offers.GET("/:id", offerHandler.GetOffer)
		offers.GET("/listing/:listing_id", offerHandler.GetListingOffers)
		offers.DELETE("/:id", offerHandler.CancelOffer)
	}

	// =========================
	// Orders
	// =========================
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, offerRepo, listingRepo, db)
	orderHandler := order.NewHandler(orderService)

	orders := api.Group("/orders")
	orders.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		orders.GET("/:id", orderHandler.GetOrder)
		orders.GET("/buyer", orderHandler.GetBuyerOrders)
		orders.GET("/farmer", orderHandler.GetFarmerOrders)
		orders.POST("/:offer_id/accept", orderHandler.AcceptOffer)
	}
	return router
}
