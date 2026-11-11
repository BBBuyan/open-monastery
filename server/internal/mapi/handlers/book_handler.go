package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	_ "example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/domain/models"
	_ "example.com/open-monastery/internal/mapi/apimodels"
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
// @Produce json,application/problem+json
// @Success 200 {object} models.Book
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 404 {object} apimodels.ProblemDetails
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(
	c *gin.Context,
) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid book id",
			err,
		)
	}

	result, err := h.service.GetBook(c.Request.Context(), id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// @Summary Get a page of books
// @Tags books
// @Produce json,application/problem+json
// @Param limit query int false "Page size (default 20, capped at 100)"
// @Param cursor query string false "Last book ID from the previous page"
// @Param q query string false "Search name or description"
// @Success 200 {object} appmodels.Paginated[models.Book]
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /books [get]
func (h *BookHandler) GetBooks(
	c *gin.Context,
) error {
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	books, err := h.service.GetBooks(c.Request.Context(), options)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, books)
	return nil
}

// @Summary Create a book
// @Tags books
// @Accept json
// @Produce json,application/problem+json
// @Param book body dtos.CreateBookRequest true "Book to create"
// @Success 201 {object} models.Book
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /books [post]
func (h *BookHandler) AddBook(
	c *gin.Context,
) error {
	var req dtos.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequestWithMsg(
			"invalid.body",
			err.Error(),
			err,
		)
	}

	result, err := h.service.CreateBook(
		c.Request.Context(),
		req.Name,
		req.Description,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// @Summary Get categories belonging to a book
// @Tags books
// @Produce json,application/problem+json
// @Param id path int true "Book ID"
// @Success 200 {array} models.Category
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /books/{id}/categories [get]
func (h *BookHandler) GetBookCategories(
	c *gin.Context,
) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid book id",
			err,
		)
	}
	result, err := h.service.GetBookCategories(c.Request.Context(), id)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}
