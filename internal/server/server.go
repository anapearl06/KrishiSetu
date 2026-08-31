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
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/offer"
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
	// Listings (frontend-aligned)
	// =========================

	listingRepo := listing.NewRepository(db)
	listingService := listing.NewService(listingRepo)
	listingHandler := listing.NewHandler(listingService)

	listingRoutes := api.Group("/listings")
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

	demandRoutes := api.Group("/demands")
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

	offerRoutes := api.Group("/offers")
	{
		// View offers for a listing
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

	// =========================
	// Orders
	// =========================

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	orderRoutes := api.Group("/orders")
	{
		orderRoutes.GET(
			"",
			middleware.JWTAuth(cfg.JWTSecret),
			orderHandler.ListOrders,
		)

		orderRoutes.GET(
			"/:id",
			middleware.JWTAuth(cfg.JWTSecret),
			orderHandler.GetOrder,
		)

		orderRoutes.POST(
			"",
			middleware.JWTAuth(cfg.JWTSecret),
			orderHandler.CreateOrder,
		)

		orderRoutes.PUT(
			"/:id/status",
			middleware.JWTAuth(cfg.JWTSecret),
			orderHandler.UpdateOrderStatus,
		)
	}

	return router
}