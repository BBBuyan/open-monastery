package controllers

import (
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/domain/models"
	_ "example.com/open-monastery/internal/api/apimodels"
	"example.com/open-monastery/internal/api/dtos"
	"example.com/open-monastery/internal/api/helpers"
	"example.com/open-monastery/internal/api/middlewares"
	"github.com/gin-gonic/gin"
)

type ReadingController struct {
	service *services.ReadingService
}

func NewReadingController(service *services.ReadingService) *ReadingController {
	return &ReadingController{service: service}
}

// GetReadings returns a page of the authenticated requester's readings at the specified monastery.
//
//	@Summary		Get a page of readings
//	@Description	Returns readings owned by the authenticated requester at the monastery.
//	@Security		BearerAuth
//	@Tags			readings
//	@Produce		json,application/problem+json
//	@Param			limit	query		int		false	"Page size (default 20, capped at 100)"
//	@Param			cursor	query		string	false	"Next cursor from the previous page"
//	@Param			id		path		int		true	"Monastery ID"
//	@Success		200		{object}	appmodels.Paginated[models.Reading]
//	@Failure		400		{object}	apimodels.ProblemDetails
//	@Failure		401		{object}	apimodels.ProblemDetails
//	@Failure		500		{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/readings [get]
func (h *ReadingController) GetReadings(c *gin.Context) error {
	var uri dtos.ReadingsURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetReadings(c.Request.Context(), *actor, options)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// AddReading creates a pending reading for the authenticated requester at the specified monastery.
//
//	@Summary		Create a reading
//	@Description	Creates a pending reading for the authenticated requester using a book from the monastery.
//	@Security		BearerAuth
//	@Tags			readings
//	@Accept			json
//	@Produce		json,application/problem+json
//	@Param			reading	body		dtos.CreateReadingRequest	true	"Reading to create"
//	@Param			id		path		int							true	"Monastery ID"
//	@Success		201		{object}	models.Reading
//	@Failure		400		{object}	apimodels.ProblemDetails
//	@Failure		401		{object}	apimodels.ProblemDetails
//	@Failure		404		{object}	apimodels.ProblemDetails
//	@Failure		500		{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/readings [post]
func (h *ReadingController) AddReading(c *gin.Context) error {
	var uri dtos.ReadingsURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	var req dtos.CreateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}

	result, err := h.service.CreateReading(
		c.Request.Context(),
		req.BookId,
		actor.Id,
		uri.MonasteryId,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// ClaimReading assigns an unassigned reading to the authenticated monk.
//
//	@Summary		Claim an unassigned reading
//	@Description	Assigns the reading to the authenticated monk if the monk belongs to the monastery.
//	@Security		BearerAuth
//	@Tags			readings
//	@Produce		application/problem+json
//	@Param			readingId	path	int	true	"Reading ID"
//	@Param			id			path	int	true	"Monastery ID"
//	@Success		204
//	@Failure		400	{object}	apimodels.ProblemDetails
//	@Failure		401	{object}	apimodels.ProblemDetails
//	@Failure		403	{object}	apimodels.ProblemDetails
//	@Failure		404	{object}	apimodels.ProblemDetails
//	@Failure		409	{object}	apimodels.ProblemDetails
//	@Failure		500	{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/readings/{readingId}/claim [post]
func (h *ReadingController) ClaimReading(c *gin.Context) error {
	var uri dtos.ReadingURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	if err := h.service.ClaimReading(c.Request.Context(), uri.MonasteryId, uri.ReadingId, actor.Id); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

// AssignReading assigns or reassigns a reading to a monk in the specified monastery.
//
//	@Summary		Assign or reassign a reading
//	@Description	Requires the Manager role for this monastery. The monk must belong to the monastery.
//	@Security		BearerAuth
//	@Tags			readings
//	@Accept			json
//	@Produce		application/problem+json
//	@Param			id			path	int							true	"Monastery ID"
//	@Param			readingId	path	int							true	"Reading ID"
//	@Param			assignment	body	dtos.AssignReadingRequest	true	"Monk to assign"
//	@Success		204
//	@Failure		400	{object}	apimodels.ProblemDetails
//	@Failure		401	{object}	apimodels.ProblemDetails
//	@Failure		403	{object}	apimodels.ProblemDetails
//	@Failure		404	{object}	apimodels.ProblemDetails
//	@Failure		409	{object}	apimodels.ProblemDetails
//	@Failure		500	{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/readings/{readingId}/assignment [put]
func (h *ReadingController) AssignReading(c *gin.Context) error {
	var uri dtos.ReadingURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	if actor.MonasteryId == nil || *actor.MonasteryId != uri.MonasteryId {
		return apperr.NewForbidden("monastery.forbidden")
	}
	var req dtos.AssignReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	if err := h.service.AssignReading(c.Request.Context(), uri.MonasteryId, uri.ReadingId, req.MonkId); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}
