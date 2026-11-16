package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/domain/models"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"example.com/open-monastery/internal/mapi/helpers"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: service}
}

// @Summary	Get categories
// @Tags		categories
// @Produce	json,application/problem+json
// @Param		id	path		int	true	"Monastery ID"
// @Success	200	{array}		models.Category
// @Failure	400	{object}	apimodels.ProblemDetails
// @Failure	500	{object}	apimodels.ProblemDetails
// @Router		/monasteries/{id}/categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	result, err := h.categoryService.GetCategories(c.Request.Context(), monasteryID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary	Create a category
// @Tags		categories
// @Accept		json
// @Produce	json,application/problem+json
// @Param		category	body		dtos.CreateCategoryRequest	true	"Category to create"
// @Param		id			path		int							true	"Monastery ID"
// @Success	201			{object}	models.Category
// @Failure	400			{object}	apimodels.ProblemDetails
// @Failure	404			{object}	apimodels.ProblemDetails
// @Failure	500			{object}	apimodels.ProblemDetails
// @Router		/monasteries/{id}/categories [post]
func (h *CategoryHandler) AddCategory(c *gin.Context) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	name := strings.TrimSpace(req.Name)
	result, err := h.categoryService.CreateCategory(c.Request.Context(), monasteryID, name)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// @Summary	Add a book to a category
// @Tags		categories
// @Accept		json
// @Param		categoryId	path	int								true	"Category ID"
// @Param		book		body	dtos.AddBookToCategoryRequest	true	"Book to add"
// @Produce	application/problem+json
// @Param		id	path	int	true	"Monastery ID"
// @Success	204	"Book added"
// @Failure	400	{object}	apimodels.ProblemDetails
// @Failure	404	{object}	apimodels.ProblemDetails
// @Failure	500	{object}	apimodels.ProblemDetails
// @Router		/monasteries/{id}/categories/{categoryId}/books [post]
func (h *CategoryHandler) AddBook(c *gin.Context) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid category id",
			err,
		)
	}
	var req dtos.AddBookToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	if err := h.categoryService.AddBook(
		c.Request.Context(),
		monasteryID,
		id,
		req.BookId,
	); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

// @Summary	Get books belonging to a category
// @Tags		categories
// @Produce	json,application/problem+json
// @Param		categoryId	path		int		true	"Category ID"
// @Param		limit		query		int		false	"Page size (default 20, capped at 100)"
// @Param		cursor		query		string	false	"Last book ID from the previous page"
// @Param		q			query		string	false	"Search name or description"
// @Param		id			path		int		true	"Monastery ID"
// @Success	200			{object}	appmodels.Paginated[models.Book]
// @Failure	400			{object}	apimodels.ProblemDetails
// @Failure	500			{object}	apimodels.ProblemDetails
// @Router		/monasteries/{id}/categories/{categoryId}/books [get]
func (h *CategoryHandler) GetBooksByCategoryId(c *gin.Context) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid category id",
			err)
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	result, err := h.categoryService.GetBooksByCategoryId(c.Request.Context(), monasteryID, id, options)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}
