package handlers

import (
	"errors"
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/domain/models"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"example.com/open-monastery/internal/mapi/helpers"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
)

type ReadingHandler struct {
	service *services.ReadingService
}

func NewReadingHandler(service *services.ReadingService) *ReadingHandler {
	return &ReadingHandler{service: service}
}

// @Summary Get a page of readings
// @Description Returns readings owned by the authenticated requester.
// @Security BearerAuth
// @Tags readings
// @Produce json,application/problem+json
// @Param limit query int false "Page size (default 20, capped at 100)"
// @Param cursor query string false "Next cursor from the previous page"
// @Success 200 {object} appmodels.Paginated[models.Reading]
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 401 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /readings [get]
func (h *ReadingHandler) GetReadings(c *gin.Context) error {
	userID, ok := middlewares.UserID(c)
	if !ok {
		return errors.New("missing authenticated user ID")
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetReadings(c.Request.Context(), userID, options)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary Create a reading
// @Description Creates a pending reading for the authenticated requester.
// @Security BearerAuth
// @Tags readings
// @Accept json
// @Produce json,application/problem+json
// @Param reading body dtos.CreateReadingRequest true "Reading to create"
// @Success 201 {object} models.Reading
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 401 {object} apimodels.ProblemDetails
// @Failure 404 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /readings [post]
func (h *ReadingHandler) AddReading(c *gin.Context) error {
	userID, ok := middlewares.UserID(c)
	if !ok {
		return errors.New("missing authenticated user ID")
	}
	var req dtos.CreateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}

	result, err := h.service.CreateReading(
		c.Request.Context(),
		req.BookId,
		userID,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}
