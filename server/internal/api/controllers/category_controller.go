package controllers

import (
	"net/http"
	"strings"

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

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(service *services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: service}
}

// GetCategories lists the categories in the specified monastery.
//
//	@Summary	Get categories
//	@Tags		categories
//	@Produce	json,application/problem+json
//	@Param		id	path		int	true	"Monastery ID"
//	@Success	200	{array}		models.Category
//	@Failure	400	{object}	apimodels.ProblemDetails
//	@Failure	500	{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id}/categories [get]
func (h *CategoryController) GetCategories(c *gin.Context) error {
	var uri dtos.CategoriesURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	result, err := h.categoryService.GetCategories(c.Request.Context(), uri.MonasteryId)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}

// AddCategory creates a category in the specified monastery.
//
//	@Summary		Create a category
//	@Description	Requires the Manager role for this monastery.
//	@Tags			categories
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json,application/problem+json
//	@Param			category	body		dtos.CreateCategoryRequest	true	"Category to create"
//	@Param			id			path		int							true	"Monastery ID"
//	@Success		201			{object}	models.Category
//	@Failure		400			{object}	apimodels.ProblemDetails
//	@Failure		401			{object}	apimodels.ProblemDetails
//	@Failure		403			{object}	apimodels.ProblemDetails
//	@Failure		404			{object}	apimodels.ProblemDetails
//	@Failure		500			{object}	apimodels.ProblemDetails
//	@Router			/categories [post]
func (h *CategoryController) AddCategory(c *gin.Context) error {
	var uri dtos.CategoriesURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	if actor.MonasteryId != nil && *actor.MonasteryId != uri.MonasteryId {
		return apperr.NewForbidden("monastery.forbidden")
	}
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	name := strings.TrimSpace(req.Name)
	result, err := h.categoryService.CreateCategory(c.Request.Context(), name, *actor)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// AddBook adds an existing book to a category.
//
//	@Summary		Add a book to a category
//	@Description	Requires the Manager role for this monastery.
//	@Tags			categories
//	@Security		BearerAuth
//	@Accept			json
//	@Param			categoryId	path	int								true	"Category ID"
//	@Param			book		body	dtos.AddBookToCategoryRequest	true	"Book to add"
//	@Produce		application/problem+json
//	@Param			id	path	int	true	"Monastery ID"
//	@Success		204	"Book added"
//	@Failure		400	{object}	apimodels.ProblemDetails
//	@Failure		401	{object}	apimodels.ProblemDetails
//	@Failure		403	{object}	apimodels.ProblemDetails
//	@Failure		404	{object}	apimodels.ProblemDetails
//	@Failure		500	{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/categories/{categoryId}/books [post]
func (h *CategoryController) AddBook(c *gin.Context) error {
	var uri dtos.CategoryURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	actor, err := middlewares.AuthenticatedActor(c)
	if err != nil {
		return err
	}
	if actor.MonasteryId != nil && *actor.MonasteryId != uri.MonasteryId {
		return apperr.NewForbidden("monastery.forbidden")
	}
	var req dtos.AddBookToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	if err := h.categoryService.AddBook(
		c.Request.Context(),
		uri.CategoryId,
		req.BookId,
		*actor,
	); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

// GetBooksByCategoryId returns a page of books in the specified category.
//
//	@Summary	Get books belonging to a category
//	@Tags		categories
//	@Produce	json,application/problem+json
//	@Param		categoryId	path		int		true	"Category ID"
//	@Param		limit		query		int		false	"Page size (default 20, capped at 100)"
//	@Param		cursor		query		string	false	"Last book ID from the previous page"
//	@Param		q			query		string	false	"Search name or description"
//	@Param		id			path		int		true	"Monastery ID"
//	@Success	200			{object}	appmodels.Paginated[models.Book]
//	@Failure	400			{object}	apimodels.ProblemDetails
//	@Failure	500			{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id}/categories/{categoryId}/books [get]
func (h *CategoryController) GetBooksByCategoryId(c *gin.Context) error {
	var uri dtos.CategoryURI
	if err := c.ShouldBindUri(&uri); err != nil {
		return apperr.NewBadRequest("invalid.uri", err)
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	result, err := h.categoryService.GetBooksByCategoryId(c.Request.Context(), uri.MonasteryId, uri.CategoryId, options)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}
