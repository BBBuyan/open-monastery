package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"example.com/open-monastery/internal/mapi/helpers"
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

// @Summary	Get a monastery
// @Param		id	path	int	true	"Monastery ID"
// @Tags		monasteries
// @Accept		json
// @Produce	json,application/problem+json
// @Success	200	{object}	models.Monastery
// @Failure	400	{object}	apimodels.ProblemDetails
// @Failure	404	{object}	apimodels.ProblemDetails
// @Router		/monasteries/{id} [get]
func (h *MonasteryHandler) GetMonastery(
	c *gin.Context,
) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid monastery id",
			err,
		)
	}

	result, err := h.service.GetMonastery(c.Request.Context(), id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary	Get a page of monasteries
// @Tags		monasteries
// @Produce	json,application/problem+json
// @Param		limit	query		int		false	"Page size (default 20, capped at 100)"
// @Param		cursor	query		string	false	"Next cursor from the previous page"
// @Param		q		query		string	false	"Search by monastery name"
// @Success	200		{object}	appmodels.Paginated[models.Monastery]
// @Failure	400		{object}	apimodels.ProblemDetails
// @Failure	500		{object}	apimodels.ProblemDetails
// @Router		/monasteries [get]
func (h *MonasteryHandler) GetMonasteries(
	c *gin.Context,
) error {
	opts, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	monasteries, err := h.service.GetMonasteries(c.Request.Context(), opts)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, monasteries)
	return nil
}

// @Summary	Create a monastery
// @Tags		monasteries
// @Accept		json
// @Produce	json,application/problem+json
// @Param		monastery	body		dtos.CreateMonasteryRequest	true	"Monastery to create"
// @Success	201			{object}	models.Monastery
// @Failure	400			{object}	apimodels.ProblemDetails
// @Failure	500			{object}	apimodels.ProblemDetails
// @Router		/monasteries [post]
func (h *MonasteryHandler) AddMonastery(
	c *gin.Context,
) error {
	var req dtos.CreateMonasteryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	result, err := h.service.CreateMonastery(c.Request.Context(), req.Name)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}
