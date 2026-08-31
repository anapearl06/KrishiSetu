package offer

import (
<<<<<<< HEAD
=======
	"errors"
>>>>>>> tmp-pr-merge
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	"github.com/raaj2493/KrishiSetu/internal/middleware"
=======
>>>>>>> tmp-pr-merge
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type Handler struct {
<<<<<<< HEAD
	service         *Service
	resolveFarmerID func(listingID uint) (uint, error)
}

func NewHandler(service *Service, resolveFarmerID func(listingID uint) (uint, error)) *Handler {
	return &Handler{
		service:         service,
		resolveFarmerID: resolveFarmerID,
	}
}

func (h *Handler) CreateOffer(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	buyerID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user identity")
		return
	}

	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	farmerID, err := h.resolveFarmerID(input.ListingID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	offer, err := h.service.CreateOffer(buyerID, farmerID, input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, ToResponse(offer))
}

func (h *Handler) GetMySentOffers(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	buyerID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user identity")
=======
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// =========================
// Request DTOs
// =========================

type CreateOfferRequest struct {
	ListingID    uint    `json:"listing_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	OfferedPrice float64 `json:"offered_price" binding:"gte=0"`
	Message      string  `json:"message"`
}

// =========================
// Create Offer
// POST /api/v1/offers
// =========================

func (h *Handler) CreateOffer(c *gin.Context) {
	buyerIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	buyerID, ok := buyerIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "invalid user id")
		return
	}

	var req CreateOfferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request_body", "invalid request body")
		return
	}

	offer, err := h.service.CreateOffer(
		buyerID,
		req.ListingID,
		req.Quantity,
		req.OfferedPrice,
		req.Message,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrListingNotFound):
			response.Error(c, http.StatusNotFound, "listing_not_found", err.Error())

		case errors.Is(err, ErrListingNotActive),
			errors.Is(err, ErrQuantityUnavailable),
			errors.Is(err, ErrInvalidQuantity),
			errors.Is(err, ErrInvalidPrice):
			response.Error(c, http.StatusBadRequest, "invalid_offer_data", err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, "server_error", "failed to create offer")
		}

		return
	}

	response.Success(c, http.StatusCreated, offer)
}

// =========================
// Get Offer
// GET /api/v1/offers/:id
// =========================

func (h *Handler) GetOffer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_offer_id", "invalid offer id")
		return
	}

	offer, err := h.service.GetOffer(uint(id))
	if err != nil {
		if errors.Is(err, ErrOfferNotFound) {
			response.Error(c, http.StatusNotFound, "offer_not_found", err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get offer")
		return
	}

	response.Success(c, http.StatusOK, offer)
}

// =========================
// Get My Offers
// GET /api/v1/offers/my
// =========================

func (h *Handler) GetMyOffers(c *gin.Context) {
	buyerIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	buyerID, ok := buyerIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
>>>>>>> tmp-pr-merge
		return
	}

	offers, err := h.service.GetBuyerOffers(buyerID)
	if err != nil {
<<<<<<< HEAD
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusOK, ToResponseList(offers))
}

func (h *Handler) GetFarmerOffers(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	farmerID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user identity")
		return
	}

	offers, err := h.service.GetFarmerOffers(farmerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusOK, ToResponseList(offers))
}

func (h *Handler) RespondOffer(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	farmerID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user identity")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid offer id")
		return
	}

	var input RespondInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	offer, err := h.service.RespondOffer(farmerID, uint(id), input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	response.Success(c, http.StatusOK, ToResponse(offer))
}

func (h *Handler) GetListingOffers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid listing id")
		return
	}

	offers, err := h.service.GetListingOffers(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusOK, ToResponseList(offers))
=======
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get offers")
		return
	}

	response.Success(c, http.StatusOK, offers)
}

// =========================
// Get Listing Offers
// GET /api/v1/offers/listing/:listing_id
// =========================

func (h *Handler) GetListingOffers(c *gin.Context) {
	listingID, err := strconv.ParseUint(
		c.Param("listing_id"),
		10,
		64,
	)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_listing_id", "invalid listing id")
		return
	}

	offers, err := h.service.GetListingOffers(uint(listingID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get listing offers")
		return
	}

	response.Success(c, http.StatusOK, offers)
}

// =========================
// Cancel Offer
// DELETE /api/v1/offers/:id
// =========================

func (h *Handler) CancelOffer(c *gin.Context) {
	buyerIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	buyerID, ok := buyerIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	offerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_offer_id", "invalid offer id")
		return
	}

	err = h.service.CancelOffer(
		uint(offerID),
		buyerID,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrOfferNotFound):
			response.Error(c, http.StatusNotFound, "offer_not_found", err.Error())

		case errors.Is(err, ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "forbidden", err.Error())

		case errors.Is(err, ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, "invalid_offer_status", err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, "server_error", "failed to cancel offer")
		}

		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "offer cancelled successfully",
	})
>>>>>>> tmp-pr-merge
}
