package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/config"
	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/listing"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/offer"
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

	// =========================
	// Demand (buyer creates demand)
	// =========================

	demandRepo := demand.NewRepository(db)
	demandService := demand.NewService(demandRepo)
	demandHandler := demand.NewHandler(demandService)

	demandRoutes := router.Group("/api/v1/demands")
	{
		// Public: browse all demands
		demandRoutes.GET(
			"",
			demandHandler.ListDemands,
		)

		// Protected: buyer demand management (JWT)
		demandRoutes.POST(
			"",
			middleware.JWTAuth(cfg.JWTSecret),
			demandHandler.CreateDemand,
		)

		demandRoutes.GET(
			"/my",
			middleware.JWTAuth(cfg.JWTSecret),
			demandHandler.GetMyDemands,
		)

		demandRoutes.GET(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			demandHandler.GetDemand,
		)

		demandRoutes.PUT(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			demandHandler.UpdateDemand,
		)

		demandRoutes.DELETE(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			demandHandler.CancelDemand,
		)
	}

	// =========================
	// Offers
	// =========================

	offerRepo := offer.NewRepository(db)
	offerService := offer.NewService(offerRepo)
	offerHandler := offer.NewHandler(offerService, func(listingID uint) (uint, error) {
		listing, err := listingService.GetListing(listingID)
		if err != nil {
			return 0, err
		}
		return listing.FarmerID, nil
	})

	offerRoutes := router.Group("/api/v1/offers")
	{
		// Public: view offers for a listing (optional auth handled in handlers)
		offerRoutes.GET(
			"/listing/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			offerHandler.GetListingOffers,
		)

		// Protected: buyer offers
		offerRoutes.POST(
			"",
			middleware.JWTAuth(cfg.JWTSecret),
			offerHandler.CreateOffer,
		)

		offerRoutes.GET(
			"/my-sent",
			middleware.JWTAuth(cfg.JWTSecret),
			offerHandler.GetMySentOffers,
		)

		offerRoutes.GET(
			"/my-received",
			middleware.JWTAuth(cfg.JWTSecret),
			offerHandler.GetFarmerOffers,
		)

		offerRoutes.PUT(
			"/:id/respond",
			middleware.JWTAuth(cfg.JWTSecret),
			offerHandler.RespondOffer,
		)
	}

	return router
}
