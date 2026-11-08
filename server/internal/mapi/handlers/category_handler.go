package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
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
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	result, err := h.categoryService.GetCategories(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary Create a category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dtos.CreateCategoryRequest true "Category to create"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) AddCategory(c *gin.Context) {
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.NewBadRequest("invalid request body", err))
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		_ = c.Error(apperr.NewBadRequest("name is required", nil))
		return
	}
	result, err := h.categoryService.CreateCategory(c.Request.Context(), name)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary Add a book to a category
// @Tags categories
// @Accept json
// @Param id path int true "Category ID"
// @Param book body dtos.AddBookToCategoryRequest true "Book to add"
// @Success 204 "Book added"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/books [post]
func (h *CategoryHandler) AddBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest("invalid category id", err))
		return
	}
	var req dtos.AddBookToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.NewBadRequest("invalid request body", err))
		return
	}
	if req.BookId <= 0 {
		_ = c.Error(apperr.NewBadRequest("book_id must be positive", nil))
		return
	}
	if err := h.categoryService.AddBook(c.Request.Context(), id, req.BookId); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get books belonging to a category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Param limit query int false "Page size (default 20, capped at 100)"
// @Param cursor query string false "Last book ID from the previous page"
// @Param q query string false "Search name or description"
// @Success 200 {array} models.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/books [get]
func (h *CategoryHandler) GetBooksByCategoryId(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest("invalid category id", err))
		return
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result, err := h.categoryService.GetBooksByCategoryId(c.Request.Context(), id, options)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}
