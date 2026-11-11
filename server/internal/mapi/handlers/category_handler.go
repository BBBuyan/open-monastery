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

// @Summary Get categories
// @Tags categories
// @Produce json,application/problem+json
// @Success 200 {array} models.Category
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) error {
	result, err := h.categoryService.GetCategories(c.Request.Context())
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary Create a category
// @Tags categories
// @Accept json
// @Produce json,application/problem+json
// @Param category body dtos.CreateCategoryRequest true "Category to create"
// @Success 201 {object} models.Category
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /categories [post]
func (h *CategoryHandler) AddCategory(c *gin.Context) error {
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	name := strings.TrimSpace(req.Name)
	result, err := h.categoryService.CreateCategory(c.Request.Context(), name)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// @Summary Add a book to a category
// @Tags categories
// @Accept json
// @Param id path int true "Category ID"
// @Param book body dtos.AddBookToCategoryRequest true "Book to add"
// @Produce application/problem+json
// @Success 204 "Book added"
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 404 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /categories/{id}/books [post]
func (h *CategoryHandler) AddBook(c *gin.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
		id,
		req.BookId,
	); err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

// @Summary Get books belonging to a category
// @Tags categories
// @Produce json,application/problem+json
// @Param id path int true "Category ID"
// @Param limit query int false "Page size (default 20, capped at 100)"
// @Param cursor query string false "Last book ID from the previous page"
// @Param q query string false "Search name or description"
// @Success 200 {object} appmodels.Paginated[models.Book]
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /categories/{id}/books [get]
func (h *CategoryHandler) GetBooksByCategoryId(c *gin.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid category id",
			err)
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	result, err := h.categoryService.GetBooksByCategoryId(c.Request.Context(), id, options)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}
