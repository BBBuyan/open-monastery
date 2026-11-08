package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/mapi/dtos"
	"example.com/open-monastery/internal/mapi/helpers"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// @Summary Get Book
// @Param id path int true "Book ID"
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest(
			"invalid book id",
			err,
		))
		return
	}

	result, err := h.service.GetBook(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Get a page of books
// @Tags books
// @Produce json
// @Param limit query int false "Page size (default 20, capped at 100)"
// @Param cursor query string false "Last book ID from the previous page"
// @Param q query string false "Search name or description"
// @Success 200 {object} appmodels.Paginated[models.Book]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [get]
func (h *BookHandler) GetBooks(
	c *gin.Context,
) {
	options, err := helpers.ToListOpts(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	books, err := h.service.GetBooks(c.Request.Context(), options)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Create a book
// @Tags books
// @Accept json
// @Produce json
// @Param book body dtos.CreateBookRequest true "Book to create"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) AddBook(
	c *gin.Context,
) {
	var req dtos.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.NewBadRequest(
			"invalid request body",
			err,
		))
		return
	}
	if req.Name == "" {
		_ = c.Error(apperr.NewBadRequest(
			"name is required",
			nil,
		))
		return
	}

	result, err := h.service.CreateBook(
		c.Request.Context(),
		req.Name,
		req.Description,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary Get categories belonging to a book
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {array} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id}/categories [get]
func (h *BookHandler) GetBookCategories(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest("invalid book id", err))
		return
	}
	result, err := h.service.GetBookCategories(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}
