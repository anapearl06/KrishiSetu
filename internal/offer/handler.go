package offer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type Handler struct {
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
		return
	}

	offers, err := h.service.GetBuyerOffers(buyerID)
	if err != nil {
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
}
