package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/models"
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
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		_ = c.Error(apperr.NewBadRequest(
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
// @Tags users
// @Accept json
// @Produce json
// @Param user body api.CreateUserRequest true "User to create"
// @Success 201 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.NewBadRequest(
			"invalid request body",
			err,
		))
		return
	}
	if req.FirstName == "" || req.LastName == "" {
		_ = c.Error(apperr.NewBadRequest(
			"first_name and last_name are required",
			nil,
		))
		return
	}

	var (
		result *models.User
		err    error
	)

	switch req.UserType {
	case models.Customer:
		result, err = h.service.CreateUser(
			c.Request.Context(),
			req.FirstName,
			req.LastName,
		)
	case models.Monk:
		if req.MonasteryId == nil || *req.MonasteryId <= 0 {
			_ = c.Error(apperr.NewBadRequest(
				"a positive monastery_id is required for monks",
				nil,
			))
			return
		}
		result, err = h.service.CreateMonk(
			c.Request.Context(),
			req.FirstName,
			req.LastName,
			*req.MonasteryId,
		)
	default:
		_ = c.Error(apperr.NewBadRequest(
			"user_type must be customer or monk",
			nil,
		))
		return
	}

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
