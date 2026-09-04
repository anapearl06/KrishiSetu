package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/market"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/offer"
	"github.com/raaj2493/KrishiSetu/internal/order"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
	"github.com/raaj2493/KrishiSetu/internal/matching"

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
		farmerRoutes.POST("/register", farmerHandler.Register)
		farmerRoutes.POST("/login", farmerHandler.Login)
		farmerRoutes.GET("/me", middleware.JWTAuth(cfg.JWTSecret), farmerHandler.Me)
		farmerRoutes.PUT("/me", middleware.JWTAuth(cfg.JWTSecret), farmerHandler.UpdateProfile)
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
		buyerRoutes.POST("/register", buyerHandler.Register)
		buyerRoutes.POST("/login", buyerHandler.Login)
		buyerRoutes.GET("/me", middleware.JWTAuth(cfg.JWTSecret), buyerHandler.Me)
		buyerRoutes.PUT("/me", middleware.JWTAuth(cfg.JWTSecret), buyerHandler.UpdateProfile)
	}

	// =========================
	// Auth (unified profile)
	// =========================

	authHandler := newAuthHandler(farmerService, buyerService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.GET("/me", middleware.JWTAuth(cfg.JWTSecret), authHandler.Me)
	}

	// =========================
	// Demands
	// =========================

	demandRepo := demand.NewRepository(db)
	demandService := demand.NewService(demandRepo)
	demandHandler := demand.NewHandler(demandService)

	demandRoutes := api.Group("/demands")
	demandRoutes.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		demandRoutes.POST("", demandHandler.CreateDemand)
		demandRoutes.GET("", demandHandler.ListDemands)
		demandRoutes.GET("/my", demandHandler.GetMyDemands)
		demandRoutes.GET("/:id", demandHandler.GetDemand)
		demandRoutes.PUT("/:id", demandHandler.UpdateDemand)
		demandRoutes.DELETE("/:id", demandHandler.CancelDemand)
	}

	// =========================
	// Listings
	// =========================

	listingRepo := listing.NewRepository(db)
	listingService := listing.NewService(listingRepo)
	listingHandler := listing.NewHandler(listingService)

	listingRoutes := api.Group("/listings")
	{
		listingRoutes.GET("", listingHandler.Browse)
		listingRoutes.POST("", middleware.JWTAuth(cfg.JWTSecret), listingHandler.Create)
		listingRoutes.GET("/my", middleware.JWTAuth(cfg.JWTSecret), listingHandler.GetMyListings)
		listingRoutes.PUT("/:id", middleware.JWTAuth(cfg.JWTSecret), listingHandler.Update)
		listingRoutes.DELETE("/:id", middleware.JWTAuth(cfg.JWTSecret), listingHandler.Delete)
	}

	// =========================
	// Offers
	// =========================
	offerRepo := offer.NewRepository(db)
	offerService := offer.NewService(offerRepo, db)
	offerHandler := offer.NewHandler(offerService)

	// =========================
	// Orders
	// =========================
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, offerRepo, listingRepo, db)
	orderHandler := order.NewHandler(orderService)

	offerRoutes := api.Group("/offers")
	offerRoutes.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		offerRoutes.POST("", offerHandler.CreateOffer)
		offerRoutes.GET("/my", offerHandler.GetMyOffers)
		offerRoutes.GET("/:id", offerHandler.GetOffer)
		offerRoutes.GET("/listing/:listing_id", offerHandler.GetListingOffers)
		offerRoutes.DELETE("/:id", offerHandler.CancelOffer)
		// Missing routes frontend expects
		offerRoutes.GET("/buyer", offerHandler.GetMyOffers)
		offerRoutes.GET("/farmer", offerHandler.GetFarmerOffers)
		offerRoutes.POST("/:id/accept", orderHandler.AcceptOffer)
		offerRoutes.POST("/:id/reject", offerHandler.RejectOffer)
	}

	// =========================
// Matching
// =========================

matchingRepo := matching.NewRepository(db)

matchingService := matching.NewService(
	matchingRepo,
	listingRepo,
	demandRepo,
)

matchingHandler := matching.NewHandler(
	matchingService,
)

matchingRoutes := api.Group("/matching")
matchingRoutes.Use(middleware.JWTAuth(cfg.JWTSecret))

{
	matchingRoutes.POST(
		"/listings/:listingID/demands/:demandID",
		matchingHandler.CreateMatch,
	)

	matchingRoutes.POST(
		"/listings/:listingID/generate",
		matchingHandler.GenerateMatchesForListing,
	)

	matchingRoutes.GET(
		"/listings/:listingID",
		matchingHandler.GetMatchesForListing,
	)

	matchingRoutes.POST(
		"/demands/:demandID/generate",
		matchingHandler.GenerateMatchesForDemand,
	)

	matchingRoutes.GET(
		"/demands/:demandID",
		matchingHandler.GetMatchesForDemand,
	)
}
	// =========================
	// Market
	// =========================
	marketRepo := market.NewRepository(db)
	marketService := market.NewService(marketRepo)
	marketHandler := market.NewHandler(marketService)

	marketRoutes := api.Group("/market")
{
	marketRoutes.GET(
		"/prices/intelligence",
		marketHandler.GetPriceIntelligence,
	)

	marketRoutes.GET(
		"/prices/regional",
		marketHandler.GetRegionalPrices,
	)
}
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
