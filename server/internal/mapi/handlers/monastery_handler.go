package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

type MonasteryHandler struct {
	service *services.MonasteryService
}

func NewMonasteryHandler(service *services.MonasteryService) *MonasteryHandler {
	return &MonasteryHandler{
		service: service,
	}
}

// @Summary Get a monastery
// @Param id path int true "Monastery ID"
// @Tags monasteries
// @Accept json
// @Produce json
// @Success 200 {object} models.Monastery
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /monasteries/{id} [get]
func (h *MonasteryHandler) GetMonastery(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest(
			"invalid monastery id",
			err,
		))
		return
	}

	result, err := h.service.GetMonastery(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Create a monastery
// @Tags monasteries
// @Accept json
// @Produce json
// @Param monastery body dtos.CreateMonasteryRequest true "Monastery to create"
// @Success 201 {object} models.Monastery
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monasteries [post]
func (h *MonasteryHandler) AddMonastery(
	c *gin.Context,
) {
	var req dtos.CreateMonasteryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	result, err := h.service.CreateMonastery(c.Request.Context(), req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
