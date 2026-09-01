package offer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type Handler struct {
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
		return
	}

	offers, err := h.service.GetBuyerOffers(buyerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get offers")
		return
	}

	c.JSON(http.StatusOK, offers)
}

func (h *Handler) GetFarmerOffers(c *gin.Context) {
	farmerIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	farmerID, ok := farmerIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	offers, err := h.service.GetFarmerOffers(farmerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get offers")
		return
	}

	c.JSON(http.StatusOK, offers)
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
}

// =========================
// Reject Offer
// POST /api/v1/offers/:id/reject
// =========================

func (h *Handler) RejectOffer(c *gin.Context) {
	farmerIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	farmerID, ok := farmerIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	offerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_offer_id", "invalid offer id")
		return
	}

	err = h.service.RejectOffer(
		uint(offerID),
		farmerID,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrOfferNotFound),
			errors.Is(err, ErrListingNotFound):
			response.Error(c, http.StatusNotFound, "offer_not_found", err.Error())

		case errors.Is(err, ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "forbidden", err.Error())

		case errors.Is(err, ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, "invalid_offer_status", err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, "server_error", "failed to reject offer")
		}

		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "offer rejected successfully",
	})
}
