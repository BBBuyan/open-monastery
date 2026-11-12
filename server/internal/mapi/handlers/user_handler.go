package handlers

import (
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	_ "example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// @Summary Create a user
// @Description Requires first_name, last_name, password, and at least one of email or phone_number. Creates a customer.
// @Tags users
// @Accept json
// @Produce json,application/problem+json
// @Param user body dtos.CreateUserRequest true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 409 {object} apimodels.ProblemDetails
// @Failure 422 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /users [post]
func (h *UserHandler) RegisterUser(c *gin.Context) error {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	result, err := h.service.CreateUser(
		c.Request.Context(), req.FirstName, req.LastName,
		req.Email, req.PhoneNumber, req.Password,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// @Summary Create a monk
// @Description Requires first_name, last_name, password, a positive monastery_id, and at least one of email or phone_number.
// @Tags monks
// @Accept json
// @Produce json,application/problem+json
// @Param monk body dtos.CreateMonkRequest true "Monk to create"
// @Success 201 {object} models.User
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 409 {object} apimodels.ProblemDetails
// @Failure 422 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /monks [post]
func (h *UserHandler) RegisterMonk(c *gin.Context) error {
	var req dtos.CreateMonkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	monk, err := h.service.CreateMonk(
		c.Request.Context(), req.FirstName, req.LastName,
		req.Email, req.PhoneNumber, req.MonasteryId, req.Password,
	)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, monk)
	return nil
}
