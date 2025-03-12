package handler

import (
	"errors"

	"github.com/EputraP/kfc_be/internal/dto"
	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/service"
	"github.com/EputraP/kfc_be/internal/util/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

type AuthHandlerConfig struct {
	AuthService service.AuthService
}

func NewAuthHandler(config AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{
		authService: config.AuthService,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var registerBody dto.RegisterBody

	if err := c.ShouldBindJSON(&registerBody); err != nil {
		response.Error(c, 400, errs.InvalidRequestBody.Error())
		return
	}

	resp, err := h.authService.CreateUser(&registerBody)

	if err != nil {
		if errors.Is(err, errs.UsernameAlreadyUsed) ||
			errors.Is(err, errs.PasswordContainUsername) {
			response.Error(c, 400, err.Error())
			return
		}

		response.UnknownError(c, err)
		return
	}

	response.JSON(c, 201, "Register Success", resp)
}
