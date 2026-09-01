package market

import (
	"net/http"

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

func (h *Handler) GetPriceIntelligence(c *gin.Context) {
	commodity := c.Query("commodity")

	if commodity == "" {
		response.Error(
			c,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"commodity is required",
		)
		return
	}

	result, err := h.service.GetPriceIntelligence(
		c.Request.Context(),
		commodity,
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"MARKET_ERROR",
			"failed to get price intelligence",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		result,
	)
}