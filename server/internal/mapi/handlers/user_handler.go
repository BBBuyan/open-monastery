package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// @Summary Get a user
// @Param id path int true "User ID"
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest(
			"err",
			"invalid user id",
			err,
		))
		return
	}

	result, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Create a user
// @Description Requires first_name, last_name, and at least one of email or phone_number. Creates a customer.
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.CreateUserRequest true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	result, err := h.service.CreateUser(
		c.Request.Context(), req.FirstName, req.LastName,
		req.Email, req.PhoneNumber,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary Create a monk
// @Description Requires first_name, last_name, a positive monastery_id, and at least one of email or phone_number.
// @Tags monks
// @Accept json
// @Produce json
// @Param monk body dtos.CreateMonkRequest true "Monk to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monks [post]
func (h *UserHandler) RegisterMonk(c *gin.Context) {
	var req dtos.CreateMonkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	monk, err := h.service.CreateMonk(
		c.Request.Context(), req.FirstName, req.LastName,
		req.Email, req.PhoneNumber, req.MonasteryId,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, monk)
}
