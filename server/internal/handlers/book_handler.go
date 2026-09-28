package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/api"
	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/book"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *book.Service
}

func NewBookHandler(service *book.Service) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// @Summary Get Book
// @Param id path int true "Book ID"
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {object} book.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
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

// @Summary Create a book
// @Tags books
// @Accept json
// @Produce json
// @Param book body api.CreateBookRequest true "Book to create"
// @Success 201 {object} book.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) AddBook(
	c *gin.Context,
) {
	var req api.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"invalid request body",
			err,
		))
		return
	}
	if req.Name == "" {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
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
