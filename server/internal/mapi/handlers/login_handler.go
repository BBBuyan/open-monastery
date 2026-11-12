package handlers

import (
	"errors"
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	service *services.LoginService
}

func NewLoginHandler(service *services.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

// @Summary Log in with an email address and password
// @Tags auth
// @Accept json
// @Produce json,application/problem+json
// @Param credentials body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} apimodels.LoginResponse
// @Failure 400 {object} apimodels.ProblemDetails
// @Failure 401 {object} apimodels.ProblemDetails
// @Failure 500 {object} apimodels.ProblemDetails
// @Router /auth/login [post]
func (h *LoginHandler) Login(c *gin.Context) error {
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
