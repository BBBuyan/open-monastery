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
	monasteryService *services.MonasteryService
	service          *services.BookService
}

func NewBookHandler(service *services.BookService, monasteryService *services.MonasteryService) *BookHandler {
	return &BookHandler{
		monasteryService: monasteryService,
		service:          service,
	}
}

// GetBook retrieves a book by ID from the specified monastery.
//
//	@Summary	Get Book
//	@Param		bookId	path	int	true	"Book ID"
//	@Tags		books
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Param		id	path		int	true	"Monastery ID"
//	@Success	200	{object}	models.Book
//	@Failure	400	{object}	apimodels.ProblemDetails
//	@Failure	404	{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id}/books/{bookId} [get]
func (h *BookHandler) GetBook(
	c *gin.Context,
) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("bookId"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid book id",
			err,
		)
	}

	result, err := h.service.GetBook(c.Request.Context(), monasteryID, id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result)
	return nil
}

// GetBooks returns a page of books from the specified monastery.
//
//	@Summary	Get a page of books
//	@Tags		books
//	@Produce	json,application/problem+json
//	@Param		limit	query		int		false	"Page size (default 20, capped at 100)"
//	@Param		cursor	query		string	false	"Last book ID from the previous page"
//	@Param		q		query		string	false	"Search name or description"
//	@Param		id		path		int		true	"Monastery ID"
//	@Success	200		{object}	appmodels.Paginated[models.Book]
//	@Failure	400		{object}	apimodels.ProblemDetails
//	@Failure	500		{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id}/books [get]
func (h *BookHandler) GetBooks(
	c *gin.Context,
) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	options, err := helpers.ToListOpts(c)
	if err != nil {
		return err
	}
	books, err := h.service.GetBooks(c.Request.Context(), monasteryID, options)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, books)
	return nil
}

// AddBook creates a book in the specified monastery.
//
//	@Summary		Create a book
//	@Description	Requires the Manager role for this monastery or the Admin role.
//	@Tags			books
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json,application/problem+json
//	@Param			book	body		dtos.CreateBookRequest	true	"Book to create"
//	@Param			id		path		int						true	"Monastery ID"
//	@Success		201		{object}	models.Book
//	@Failure		400		{object}	apimodels.ProblemDetails
//	@Failure		401		{object}	apimodels.ProblemDetails
//	@Failure		403		{object}	apimodels.ProblemDetails
//	@Failure		404		{object}	apimodels.ProblemDetails
//	@Failure		500		{object}	apimodels.ProblemDetails
//	@Router			/monasteries/{id}/books [post]
func (h *BookHandler) AddBook(
	c *gin.Context,
) error {
	monasteryID, err := monasteryForWrite(c, h.monasteryService)
	if err != nil {
		return err
	}
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
		monasteryID,
		req.Name,
		req.Description,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// GetBookCategories returns the categories associated with a book.
//
//	@Summary	Get categories belonging to a book
//	@Tags		books
//	@Produce	json,application/problem+json
//	@Param		bookId	path		int	true	"Book ID"
//	@Param		id		path		int	true	"Monastery ID"
//	@Success	200		{array}		models.Category
//	@Failure	400		{object}	apimodels.ProblemDetails
//	@Failure	500		{object}	apimodels.ProblemDetails
//	@Router		/monasteries/{id}/books/{bookId}/categories [get]
func (h *BookHandler) GetBookCategories(
	c *gin.Context,
) error {
	monasteryID, err := helpers.MonasteryID(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("bookId"), 10, 64)
	if err != nil || id <= 0 {
		return apperr.NewBadRequest(
			"invalid book id",
			err,
		)
	}
	result, err := h.service.GetBookCategories(c.Request.Context(), monasteryID, id)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}
