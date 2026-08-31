package order

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

// GET /api/v1/orders/:id
func (h *Handler) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_order_id", "invalid order id")
		return
	}

	order, err := h.service.GetOrder(uint(orderID))
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, "order_not_found", err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get order")
		return
	}

	response.Success(c, http.StatusOK, order)
}

// GET /api/v1/orders/buyer
func (h *Handler) GetBuyerOrders(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	buyerID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	orders, err := h.service.GetBuyerOrders(buyerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get buyer orders")
		return
	}

	response.Success(c, http.StatusOK, orders)
}

// GET /api/v1/orders/farmer
func (h *Handler) GetFarmerOrders(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	farmerID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	orders, err := h.service.GetFarmerOrders(farmerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", "failed to get farmer orders")
		return
	}

	response.Success(c, http.StatusOK, orders)
}

// POST /api/v1/orders/:offer_id/accept
func (h *Handler) AcceptOffer(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized_access", "unauthorized")
		return
	}

	farmerID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid_user_id", "invalid user id")
		return
	}

	offerID, err := strconv.ParseUint(c.Param("offer_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_offer_id", "invalid offer id")
		return
	}

	order, err := h.service.AcceptOffer(
		c.Request.Context(),
		uint(offerID),
		farmerID,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrOfferNotFound),
			errors.Is(err, ErrListingNotFound):
			response.Error(c, http.StatusNotFound, "not_found", err.Error())

		case errors.Is(err, ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "forbidden", err.Error())

		case errors.Is(err, ErrOfferNotPending),
			errors.Is(err, ErrListingNotActive),
			errors.Is(err, ErrQuantityUnavailable):
			response.Error(c, http.StatusBadRequest, "invalid_offer_state", err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, "server_error", "failed to accept offer")
		}

		return
	}

	response.Success(c, http.StatusCreated, order)
}
