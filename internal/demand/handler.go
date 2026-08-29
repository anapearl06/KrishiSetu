package demand

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateDemandRequest struct {
	CropName    string  `json:"crop_name" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"required"`
	TargetPrice float64 `json:"target_price" binding:"required,gt=0"`
	State       string  `json:"state" binding:"required"`
	District    string  `json:"district" binding:"required"`
	RequiredBy  string  `json:"required_by" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateDemand(c *gin.Context) {
	var req CreateDemandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	buyerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	buyerID, ok := buyerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	requiredBy, err := parseDate(req.RequiredBy)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid required_by, use YYYY-MM-DD",
		})
		return
	}

	demand := &Demand{
		CropName:    req.CropName,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		TargetPrice: req.TargetPrice,
		State:       req.State,
		District:    req.District,
		RequiredBy:  requiredBy,
	}

	if err := h.service.CreateDemand(
		c.Request.Context(),
		buyerID,
		demand,
	); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, demand)
}

func (h *Handler) GetDemand(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid demand id",
		})
		return
	}

	demand, err := h.service.GetDemand(
		c.Request.Context(),
		uint(id),
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, demand)
}

func (h *Handler) GetMyDemands(c *gin.Context) {
	buyerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	buyerID, ok := buyerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	demands, err := h.service.GetBuyerDemands(
		c.Request.Context(),
		buyerID,
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, demands)
}

func (h *Handler) UpdateDemand(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid demand id",
		})
		return
	}

	var req CreateDemandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	buyerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	buyerID, ok := buyerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	requiredBy, err := parseDate(req.RequiredBy)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid required_by, use YYYY-MM-DD",
		})
		return
	}

	demand := &Demand{
		ID:          uint(id),
		CropName:    req.CropName,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		TargetPrice: req.TargetPrice,
		State:       req.State,
		District:    req.District,
		RequiredBy:  requiredBy,
	}

	if err := h.service.UpdateDemand(
		c.Request.Context(),
		buyerID,
		demand,
	); err != nil {
		handleServiceError(c, err)
		return
	}

	updatedDemand, err := h.service.GetDemand(
		c.Request.Context(),
		uint(id),
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedDemand)
}

func (h *Handler) CancelDemand(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid demand id",
		})
		return
	}

	buyerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	buyerID, ok := buyerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	if err := h.service.CancelDemand(
		c.Request.Context(),
		buyerID,
		uint(id),
	); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "demand cancelled successfully",
	})
}

func (h *Handler) ListDemands(c *gin.Context) {
	filters := DemandFilters{
		CropName: c.Query("crop_name"),
		State:    c.Query("state"),
		District: c.Query("district"),
		Status:   c.Query("status"),
	}

	if value := c.Query("min_price"); value != "" {
		price, err := strconv.ParseFloat(value, 64)

		if err != nil || price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid min_price",
			})
			return
		}

		filters.MinPrice = &price
	}

	if value := c.Query("max_price"); value != "" {
		price, err := strconv.ParseFloat(value, 64)

		if err != nil || price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid max_price",
			})
			return
		}

		filters.MaxPrice = &price
	}

	demands, err := h.service.ListDemands(
		c.Request.Context(),
		filters,
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, demands)
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrDemandNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrDemandNotActive):
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrInvalidQuantity),
		errors.Is(err, ErrInvalidPrice),
		errors.Is(err, ErrInvalidCropName),
		errors.Is(err, ErrInvalidUnit),
		errors.Is(err, ErrInvalidLocation),
		errors.Is(err, ErrInvalidRequiredBy):

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
