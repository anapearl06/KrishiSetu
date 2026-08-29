package listing

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateListingRequest struct {
	CropName      string  `json:"crop_name" binding:"required"`
	Quantity      float64 `json:"quantity" binding:"required,gt=0"`
	Unit          string  `json:"unit" binding:"required"`
	ExpectedPrice float64 `json:"expected_price" binding:"required,gt=0"`
	QualityGrade  string  `json:"quality_grade"`
	State         string  `json:"state" binding:"required"`
	District      string  `json:"district" binding:"required"`
	HarvestDate   string  `json:"harvest_date" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateListing(c *gin.Context) {
	var req CreateListingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	farmerIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	farmerID, ok := farmerIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	harvestDate, err := parseHarvestDate(req.HarvestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid harvest_date, use YYYY-MM-DD",
		})
		return
	}

	listing := &CropListing{
		CropName:      req.CropName,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		ExpectedPrice: req.ExpectedPrice,
		QualityGrade:  req.QualityGrade,
		State:         req.State,
		District:      req.District,
		HarvestDate:   harvestDate,
	}

	if err := h.service.CreateListing(
		c.Request.Context(),
		farmerID,
		listing,
	); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, listing)
}

func (h *Handler) GetListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid listing id",
		})
		return
	}

	listing, err := h.service.GetListing(
		c.Request.Context(),
		uint(id),
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, listing)
}

func (h *Handler) GetMyListings(c *gin.Context) {
	farmerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	farmerID, ok := farmerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	listings, err := h.service.GetFarmerListings(
		c.Request.Context(),
		farmerID,
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, listings)
}

func (h *Handler) ListListings(c *gin.Context) {
	filters := ListingFilters{
		CropName: c.Query("crop"),
		State:    c.Query("state"),
		District: c.Query("district"),
		Status:   c.Query("status"),
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &val
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &val
		}
	}

	listings, err := h.service.ListListings(c.Request.Context(), filters)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, listings)
}

func (h *Handler) UpdateListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid listing id",
		})
		return
	}

	var req CreateListingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	farmerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	farmerID, ok := farmerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	harvestDate, err := parseHarvestDate(req.HarvestDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid harvest_date, use YYYY-MM-DD",
		})
		return
	}

	listing := &CropListing{
		ID:            uint(id),
		CropName:      req.CropName,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		ExpectedPrice: req.ExpectedPrice,
		QualityGrade:  req.QualityGrade,
		State:         req.State,
		District:      req.District,
		HarvestDate:   harvestDate,
	}

	if err := h.service.UpdateListing(
		c.Request.Context(),
		farmerID,
		listing,
	); err != nil {
		handleServiceError(c, err)
		return
	}

	updatedListing, err := h.service.GetListing(
		c.Request.Context(),
		uint(id),
	)

	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedListing)
}

func (h *Handler) CancelListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid listing id",
		})
		return
	}

	farmerIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	farmerID, ok := farmerIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identity",
		})
		return
	}

	if err := h.service.CancelListing(
		c.Request.Context(),
		farmerID,
		uint(id),
	); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "listing cancelled successfully",
	})
}

func parseHarvestDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrListingNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrListingNotActive):
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, ErrInvalidQuantity),
		errors.Is(err, ErrInvalidPrice),
		errors.Is(err, ErrInvalidCropName),
		errors.Is(err, ErrInvalidUnit),
		errors.Is(err, ErrInvalidLocation),
		errors.Is(err, ErrInvalidHarvestDate):

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}