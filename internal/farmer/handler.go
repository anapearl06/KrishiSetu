package farmer

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

	farmer, err := h.service.Register(input)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		farmer,
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

	result, err := h.service.Login(input)
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
			"farmer": gin.H{
				"id":       result.Farmer.ID,
				"name":     result.Farmer.Name,
				"phone":    result.Farmer.Phone,
				"state":    result.Farmer.State,
				"district": result.Farmer.District,
			},
		},
	)
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"user_id": userID,
			"role":    c.MustGet("role"),
		},
	)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
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

	farmer, err := h.service.UpdateProfile(id, input)
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
		farmer,
	)
}
