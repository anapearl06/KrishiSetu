package matching

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raaj2493/KrishiSetu/internal/middleware"
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

// currentUser extracts the authenticated user id and role from the
// Gin context, following the project's JWTAuth middleware conventions.
// Returns the user id and role, or an error if the identity is missing.
func currentUser(c *gin.Context) (uint, string, bool) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return 0, "", false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, "", false
	}

	roleValue, _ := c.Get(middleware.RoleKey)
	role, _ := roleValue.(string)

	return userID, role, true
}

// CreateMatch creates a match between a listing and a demand.
func (h *Handler) CreateMatch(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"missing or invalid user identity",
		)
		return
	}

	listingID, err := parseID(c.Param("listingID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid listing id",
		)
		return
	}

	demandID, err := parseID(c.Param("demandID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid demand id",
		)
		return
	}

	match, err := h.service.CreateMatch(
		c.Request.Context(),
		userID,
		role,
		listingID,
		demandID,
	)
	if err != nil {
		handleMatchingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		ToResponse(match),
	)
}

// GenerateMatchesForListing generates relevant buyer matches
// for a farmer listing.
func (h *Handler) GenerateMatchesForListing(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"missing or invalid user identity",
		)
		return
	}

	listingID, err := parseID(c.Param("listingID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid listing id",
		)
		return
	}

	matches, err := h.service.GenerateMatchesForListing(
		c.Request.Context(),
		userID,
		listingID,
	)
	if err != nil {
		handleMatchingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		ToResponseList(matches),
	)
}

// GenerateMatchesForDemand generates relevant farmer matches
// for a buyer demand.
func (h *Handler) GenerateMatchesForDemand(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"missing or invalid user identity",
		)
		return
	}

	demandID, err := parseID(c.Param("demandID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid demand id",
		)
		return
	}

	matches, err := h.service.GenerateMatchesForDemand(
		c.Request.Context(),
		userID,
		demandID,
	)
	if err != nil {
		handleMatchingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		ToResponseList(matches),
	)
}

// GetMatchesForListing returns matches for a farmer listing.
func (h *Handler) GetMatchesForListing(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"missing or invalid user identity",
		)
		return
	}

	listingID, err := parseID(c.Param("listingID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid listing id",
		)
		return
	}

	matches, err := h.service.GetMatchesForListing(
		c.Request.Context(),
		userID,
		listingID,
	)
	if err != nil {
		handleMatchingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		ToResponseList(matches),
	)
}

// GetMatchesForDemand returns matches for a buyer demand.
func (h *Handler) GetMatchesForDemand(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"missing or invalid user identity",
		)
		return
	}

	demandID, err := parseID(c.Param("demandID"))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"BAD_REQUEST",
			"invalid demand id",
		)
		return
	}

	matches, err := h.service.GetMatchesForDemand(
		c.Request.Context(),
		userID,
		demandID,
	)
	if err != nil {
		handleMatchingError(c, err)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		ToResponseList(matches),
	)
}

func parseID(value string) (uint, error) {
	id, err := strconv.ParseUint(
		value,
		10,
		32,
	)

	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}

	return uint(id), nil
}

func handleMatchingError(
	c *gin.Context,
	err error,
) {
	switch {
	case errors.Is(err, ErrListingNotFound):
		response.Error(
			c,
			http.StatusNotFound,
			"LISTING_NOT_FOUND",
			err.Error(),
		)

	case errors.Is(err, ErrDemandNotFound):
		response.Error(
			c,
			http.StatusNotFound,
			"DEMAND_NOT_FOUND",
			err.Error(),
		)

	case errors.Is(err, ErrUnauthorized):
		response.Error(
			c,
			http.StatusForbidden,
			"FORBIDDEN",
			err.Error(),
		)

	default:
		response.Error(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
	}
}
