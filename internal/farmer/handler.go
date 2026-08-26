package farmer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	State    string `json:"state"`
	District string `json:"district"`
}


type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}


func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
		)
		return
	}

	f, err := h.service.Register(RegisterInput{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
		State:    req.State,
		District: req.District,
	})

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"REGISTRATION_FAILED",
			"Failed to register farmer",
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		f,
	)
}