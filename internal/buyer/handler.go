package buyer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/middleware"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid request body",
		)
		return
	}

	buyer, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		buyer,
	)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid request body",
		)
		return
	}

	result, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"invalid credentials",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"token": result.Token,
			"buyer": gin.H{
				"id":            result.Buyer.ID,
				"name":          result.Buyer.Name,
				"phone":         result.Buyer.Phone,
				"business_name": result.Buyer.BusinessName,
				"business_type": result.Buyer.BusinessType,
				"state":         result.Buyer.State,
				"district":      result.Buyer.District,
			},
		},
	)
}

func (h *Handler) Me(c *gin.Context) {
	role, exists := c.Get(middleware.RoleKey)
	if !exists || role != "buyer" {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	id, ok := userID.(uint)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"invalid user identity",
		)
		return
	}

	buyer, err := h.service.GetProfile(
		c.Request.Context(),
		id,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusNotFound,
			"NOT_FOUND",
			"buyer not found",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		buyer,
	)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	role, exists := c.Get(middleware.RoleKey)
	if !exists || role != "buyer" {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	id, ok := userID.(uint)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"invalid user identity",
		)
		return
	}

	var input UpdateProfileInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid request body",
		)
		return
	}

	buyer, err := h.service.UpdateProfile(
		c.Request.Context(),
		id,
		input,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		buyer,
	)
}
