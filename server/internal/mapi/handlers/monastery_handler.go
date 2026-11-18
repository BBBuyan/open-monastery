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
	"example.com/open-monastery/internal/mapi/middlewares"
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

// GrantManager grants a user manager access to the specified monastery.
//
//	@Summary		Grant a user monastery manager access
//	@Description	Requires a platform admin. Appoints a customer as manager of this monastery.
//	@Tags			monasteries
//	@Security		BearerAuth
//	@Param			id		path	int	true	"Monastery ID"
//	@Param			userId	path	int	true	"User ID"
//	@Produce		application/problem+json
//	@Success		204
//	@Failure		400	{object}	apimodels.ProblemDetails
//	@Failure		401	{object}	apimodels.ProblemDetails
//	@Failure		403	{object}	apimodels.ProblemDetails
//	@Failure		404	{object}	apimodels.ProblemDetails
//	@Failure		409	{object}	apimodels.ProblemDetails
//	@Failure		500	{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/managers/{userId} [put]
func (h *MonasteryHandler) GrantManager(c *gin.Context) error {
	monasteryID, userID, err := managerPathIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.GrantManager(c.Request.Context(), monasteryID, userID); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

// RevokeManager revokes a user's manager access to the specified monastery.
//
//	@Summary		Revoke a user's monastery manager access
//	@Description	Requires a platform admin. Returns an appointed manager to the customer role.
//	@Tags			monasteries
//	@Security		BearerAuth
//	@Param			id		path	int	true	"Monastery ID"
//	@Param			userId	path	int	true	"User ID"
//	@Produce		application/problem+json
//	@Success		204
//	@Failure		400	{object}	apimodels.ProblemDetails
//	@Failure		401	{object}	apimodels.ProblemDetails
//	@Failure		403	{object}	apimodels.ProblemDetails
//	@Failure		500	{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/managers/{userId} [delete]
func (h *MonasteryHandler) RevokeManager(c *gin.Context) error {
	monasteryID, userID, err := managerPathIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.RevokeManager(c.Request.Context(), monasteryID, userID); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

func managerPathIDs(c *gin.Context) (monasteryID, userID int64, err error) {
	monasteryID, err = helpers.MonasteryID(c)
	if err != nil {
		return 0, 0, err
	}
	userID, err = strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID <= 0 {
		return 0, 0, apperr.NewBadRequest("invalid user id", err)
	}
	return monasteryID, userID, nil
}

func monasteryForWrite(c *gin.Context, service *services.MonasteryService) (int64, error) {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return 0, err
	}
	userID, ok := middlewares.UserID(c)
	if !ok {
		return 0, apperr.NewUnauthorized("authentication.required", nil)
	}
	if err := service.CheckWriteAccess(c.Request.Context(), monasteryID, userID); err != nil {
		return 0, err
	}
	return monasteryID, nil
}
