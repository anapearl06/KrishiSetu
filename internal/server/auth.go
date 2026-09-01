package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raaj2493/KrishiSetu/internal/buyer"
	"github.com/raaj2493/KrishiSetu/internal/farmer"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type authHandler struct {
	farmerService *farmer.Service
	buyerService  *buyer.Service
}

func newAuthHandler(
	farmerService *farmer.Service,
	buyerService *buyer.Service,
) *authHandler {
	return &authHandler{
		farmerService: farmerService,
		buyerService:  buyerService,
	}
}

func (h *authHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user identity")
		return
	}

	role := c.GetString(middleware.RoleKey)

	switch role {
	case "farmer":
		profile, err := h.farmerService.GetProfile(userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "failed to load profile")
			return
		}
		c.JSON(http.StatusOK, profile)

	case "buyer":
		profile, err := h.buyerService.GetProfile(c.Request.Context(), userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "failed to load profile")
			return
		}
		c.JSON(http.StatusOK, profile)

	default:
		if role == "" {
			role = "unknown"
		}
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "unsupported role: "+role)
	}
}