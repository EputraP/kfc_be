package handler

import (
	"errors"
	"net/http"

	"github.com/EputraP/kfc_be/internal/dto"
	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/service"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"github.com/EputraP/kfc_be/internal/util/response"
	"github.com/EputraP/kfc_be/internal/util/tokenprovider"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService   service.AuthService
	tokenProvider tokenprovider.JWTTokenProvider
}

type AuthHandlerConfig struct {
	AuthService   service.AuthService
	TokenProvider tokenprovider.JWTTokenProvider
}

func NewAuthHandler(config AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{
		authService:   config.AuthService,
		tokenProvider: config.TokenProvider,
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
		logger.Error("AuthHandler CreateUser", "Failed to create user", map[string]string{
			"error": err.Error(),
		})

		response.UnknownError(c, err)
		return
	}

	response.JSON(c, 201, "Register Success", resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginBody dto.LoginBody

	if err := c.ShouldBindJSON(&loginBody); err != nil {
		response.Error(c, 400, errs.InvalidRequestBody.Error())
		return
	}

	resp, err := h.authService.Login(&loginBody)
	if err != nil {
		if errors.Is(err, errs.PasswordDoesntMatch) ||
			errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, 401, errs.UsernamePasswordIncorrect.Error())
			return
		}
		logger.Error("AuthHandler Login", "Failed to login", map[string]string{
			"error": err.Error(),
		})

		response.UnknownError(c, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh-token", resp.RefreshToken, 3600*24*30, "", "", true, true)
	c.SetCookie("access-token", resp.AccesToken, 3600*24*30, "", "/", true, true)

	response.JSON(c, 200, "Login success", resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("refresh-token", "", -1, "", "", true, true)
	c.SetCookie("access-token", "", -1, "", "/", true, true)

	response.JSON(c, 200, "Logout success", nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	refreshToken, err := h.tokenProvider.ExtractToken(authHeader)
	if err != nil {
		logger.Error("AuthHandler Refresh", "Failed to extract refresh token", map[string]string{
			"error": err.Error(),
		})
		response.Error(c, 400, err.Error())
		return
	}

	token, err := h.tokenProvider.RenewAccessToken(refreshToken)
	if err != nil {
		logger.Error("AuthHandler Refresh", "Failed to renew access token", map[string]string{
			"error": err.Error(),
		})
		response.Error(c, 400, err.Error())
		return
	}

	response.JSON(c, 200, "Renew Access Token Success", *token)
}
