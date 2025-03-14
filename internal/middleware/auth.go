package middleware

import (
	"errors"
	"net/http"

	"github.com/EputraP/kfc_be/internal/constant"
	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/util/response"
	"github.com/EputraP/kfc_be/internal/util/tokenprovider"
	"github.com/gin-gonic/gin"
)

func CreateAuth(tokenChecker tokenprovider.JWTTokenProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		tokenStr, err := tokenChecker.ExtractToken(authHeader)
		if errors.Is(err, errs.InvalidBearerFormat) {
			response.Error(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		claims, err := tokenChecker.ValidateToken(tokenStr)
		if errors.Is(err, errs.InvalidToken) || errors.Is(err, errs.InvalidIssuer) {
			response.Error(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		if err != nil {
			response.UnknownError(ctx, err)
			return
		}

		ctx.Set(constant.ContextKeyUser, claims.UserClaims)
		ctx.Next()
	}
}
