package controllers

import (
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"example.com/open-monastery/internal/mapi/helpers"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
)

type MonasteryController struct {
	service *services.MonasteryService
}

func NewMonasteryController(service *services.MonasteryService) *MonasteryController {
	return &MonasteryController{
		service: service,
	}
}

// GetMonastery retrieves a monastery by ID.
//
//	@Summary	Get a monastery
//	@Param		id	path	int	true	"Monastery ID"
//	@Tags		monasteries
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Success	200	{object}	models.Monastery
//	@Failure	400	{object}	apimodels.ProblemDetails
//	@Failure	404	{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id} [get]
func (h *MonasteryController) GetMonastery(
	c *gin.Context,
) error {
	var uri dtos.MonasteriesURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}

	result, err := h.service.GetMonastery(c.Request.Context(), uri.MonasteryId)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// GetMonasteries returns a page of monasteries.
//
//	@Summary	Get a page of monasteries
//	@Tags		monasteries
//	@Produce	json,application/problem+json
//	@Param		limit	query		int		false	"Page size (default 20, capped at 100)"
//	@Param		cursor	query		string	false	"Next cursor from the previous page"
//	@Param		q		query		string	false	"Search by monastery name"
//	@Success	200		{object}	appmodels.Paginated[models.Monastery]
//	@Failure	400		{object}	apimodels.ProblemDetails
//	@Failure	500		{object}	apimodels.ProblemDetails
//	@Router		/monasteries [get]
func (h *MonasteryController) GetMonasteries(
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

// AddMonastery creates a monastery.
//
//	@Summary	Create a monastery
//	@Tags		monasteries
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Param		monastery	body		dtos.CreateMonasteryRequest	true	"Monastery to create"
//	@Success	201			{object}	models.Monastery
//	@Failure	400			{object}	apimodels.ProblemDetails
//	@Failure	401			{object}	apimodels.ProblemDetails
//	@Failure	403			{object}	apimodels.ProblemDetails
//	@Failure	500			{object}	apimodels.ProblemDetails
//	@Router		/monasteries [post]
func (h *MonasteryController) AddMonastery(
	c *gin.Context,
) error {
	var req dtos.CreateMonasteryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	result, err := h.service.CreateMonastery(c.Request.Context(), req.Name, *actor)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}
