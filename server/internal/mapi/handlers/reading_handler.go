package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

type ReadingHandler struct {
	service *services.ReadingService
}

func NewReadingHandler(service *services.ReadingService) *ReadingHandler {
	return &ReadingHandler{service: service}
}

// @Summary Get a reading
// @Param id path int true "Reading ID"
// @Tags readings
// @Accept json
// @Produce json,application/problem+json
// @Success 200 {object} models.Reading
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 404 {object} apimodels.ProblemDetails
// @Router /readings/{id} [get]
func (h *ReadingHandler) GetReading(c *gin.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequestWithMsg(
			"invalid.id",
			err.Error(),
			err,
		)
	}

	result, err := h.service.GetReading(c.Request.Context(), id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary Create a reading
// @Tags readings
// @Accept json
// @Produce json,application/problem+json
// @Param reading body dtos.CreateReadingRequest true "Reading to create"
// @Success 201 {object} models.Reading
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /readings [post]
func (h *ReadingHandler) AddReading(c *gin.Context) error {
	var req dtos.CreateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}

	result, err := h.service.CreateReading(
		c.Request.Context(),
		req.BookId,
		req.CustomerId,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}
