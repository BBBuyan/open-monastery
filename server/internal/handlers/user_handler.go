package handlers

import (
	"net/http"
	"strconv"

	"example.com/open-monastery/internal/api"
	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
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
		_ = c.Error(apperror.New(
			apperror.BadRequest,
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
	var req api.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"invalid request body",
			err,
		))
		return
	}
	if req.FirstName == "" || req.LastName == "" {
		_ = c.Error(apperror.New(
			apperror.BadRequest,
			"first_name and last_name are required",
			nil,
		))
		return
	}

	var (
		result *user.User
		err    error
	)

	switch req.UserType {
	case user.Customer:
		result, err = h.service.CreateUser(
			c.Request.Context(),
			req.FirstName,
			req.LastName,
		)
	case user.Monk:
		if req.MonasteryId == nil || *req.MonasteryId <= 0 {
			_ = c.Error(apperror.New(
				apperror.BadRequest,
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
		_ = c.Error(apperror.New(
			apperror.BadRequest,
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
