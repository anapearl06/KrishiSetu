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

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body" , err.Error())
		return
	}

	result, err := h.service.Login(LoginInput{
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		response.Error(c, 401, "invalid credentials" , err.Error())
		return
	}

	response.Success(c, 200, gin.H{
		"token": result.Token,
		"farmer": gin.H{
			"id":       result.Farmer.ID,
			"name":     result.Farmer.Name,
			"phone":    result.Farmer.Phone,
			"state":    result.Farmer.State,
			"district": result.Farmer.District,
		},
	})
}