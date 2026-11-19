package controllers

import (
	"errors"
	"net/http"

	"example.com/open-monastery/internal/api/apimodels"
	"example.com/open-monastery/internal/api/dtos"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/commands"
	"example.com/open-monastery/internal/app/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Login authenticates a user with an email address and password.
//
//	@Summary	Log in with an email address and password
//	@Tags		auth
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Param		credentials	body		dtos.LoginRequest	true	"Login credentials"
//	@Success	200			{object}	apimodels.LoginResponse
//	@Failure	400			{object}	apimodels.ProblemDetails
//	@Failure	401			{object}	apimodels.ProblemDetails
//	@Failure	500			{object}	apimodels.ProblemDetails
//	@Router		/auth/login [post]
func (h *AuthController) Login(c *gin.Context) error {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	accessToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if errors.Is(err, services.ErrInvalidCredentials) {
		return apperr.NewUnauthorized("invalid.credentials", err)
	}
	if err != nil {
		return err
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, apimodels.LoginResponse{AccessToken: accessToken, TokenType: "Bearer"})
	return nil
}

// RegisterUser registers a new user.
//
//	@Summary	Create a user
//	@Tags		auth
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Param		user	body		dtos.CreateUserRequest	true	"User to create"
//	@Success	201		{object}	models.User
//	@Failure	400		{object}	apimodels.ProblemDetails
//	@Failure	409		{object}	apimodels.ProblemDetails
//	@Failure	422		{object}	apimodels.ProblemDetails
//	@Failure	500		{object}	apimodels.ProblemDetails
//	@Router		/auth/register [post]
func (h *AuthController) RegisterUser(c *gin.Context) error {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	cmd := commands.CreateUserCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.PhoneNumber,
		Password:  req.Password,
	}
	result, err := h.service.CreateUser(c.Request.Context(), cmd)
	if err != nil {
		return err
	}
	c.JSON(http.StatusCreated, result)
	return nil
}

// RegisterMonk registers a new monk.
//
//	@Summary	Create a monk
//	@Tags		auth
//	@Accept		json
//	@Produce	json,application/problem+json
//	@Param		monk	body		dtos.CreateMonkRequest	true	"Monk to create"
//	@Success	201		{object}	models.User
//	@Failure	400		{object}	apimodels.ProblemDetails
//	@Failure	409		{object}	apimodels.ProblemDetails
//	@Failure	422		{object}	apimodels.ProblemDetails
//	@Failure	500		{object}	apimodels.ProblemDetails
//	@Router		/auth/register/monks [post]
func (h *AuthController) PromoteToMonk(c *gin.Context) error {
	var req dtos.CreateMonkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperr.NewBadRequest("invalid.body", err)
	}
	return nil
}
