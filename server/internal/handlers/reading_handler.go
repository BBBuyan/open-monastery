package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/api"
	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/reading"
	"github.com/gin-gonic/gin"
)

type ReadingHandler struct {
	service *reading.Service
}

func NewReadingHandler(service *reading.Service) *ReadingHandler {
	return &ReadingHandler{service: service}
}

// @Summary Get a reading
// @Param id path int true "Reading ID"
// @Tags readings
// @Accept json
// @Produce json
// @Success 200 {object} reading.Reading
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /readings/{id} [get]
func (h *ReadingHandler) GetReading(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"invalid reading id",
			err,
		))
		return
	}

	result, err := h.service.GetReading(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Create a reading
// @Tags readings
// @Accept json
// @Produce json
// @Param reading body api.CreateReadingRequest true "Reading to create"
// @Success 201 {object} reading.Reading
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /readings [post]
func (h *ReadingHandler) AddReading(c *gin.Context) {
	var req api.CreateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"invalid request body",
			err,
		))
		return
	}
	if req.BookId <= 0 || req.MonkId <= 0 || req.CustomerId <= 0 {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"book_id, monk_id, and customer_id must be positive",
			nil,
		))
		return
	}

	result, err := h.service.CreateReading(
		c.Request.Context(),
		req.MonkId,
		req.CustomerId,
		req.BookId,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
