package tokenprovider

import (
	"os"
	"strconv"

	"github.com/EputraP/kfc_be/internal/constant"
)

func GetProvider() JWTTokenProvider {
	issuer := constant.Issuer
	secret := os.Getenv(constant.EnvKeyJWTSecret)
	refreshTokenDurationString := os.Getenv(constant.EnvKeyRefreshTokenDuration)
	accessTokenDurationString := os.Getenv(constant.EnvKeyAccessTokenDuration)

	refreshTokenDuration, _ := strconv.Atoi(refreshTokenDurationString)
	accessTokenDuration, _ := strconv.Atoi(accessTokenDurationString)

	jwtProvider := NewJWT(issuer, secret, refreshTokenDuration, accessTokenDuration)
	return jwtProvider
}
