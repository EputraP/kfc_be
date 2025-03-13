package tokenprovider

import (
	"log"
	"strings"
	"time"

	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTTokenProvider interface {
	GenerateRefreshToken(user model.User) (string, error)
	GenerateAccessToken(user model.User) (string, error)
	ValidateToken(token string) (*JwtClaims, error)
	ExtractToken(authHeader string) (string, error)
	RenewAccessToken(refreshTokenString string) (*string, error)
}

type jwtTokenProvider struct {
	issuer               string
	secret               string
	refreshTokenDuration int
	accessTokenDuration  int
}

func NewJWT(issuer string, secret string, refreshTokenDuration int, accessTokenDuration int) JWTTokenProvider {
	return &jwtTokenProvider{
		issuer:               issuer,
		secret:               secret,
		refreshTokenDuration: refreshTokenDuration,
		accessTokenDuration:  accessTokenDuration,
	}
}

func (p *jwtTokenProvider) GenerateAccessToken(user model.User) (string, error) {
	return p.generateToken(user, time.Duration(p.accessTokenDuration)*time.Minute)
}

func (p *jwtTokenProvider) GenerateRefreshToken(user model.User) (string, error) {
	return p.generateToken(user, time.Duration(p.refreshTokenDuration)*time.Minute)
}

func (p *jwtTokenProvider) generateToken(user model.User, expiresIn time.Duration) (string, error) {
	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserClaims: UserClaims{
			UserID:   user.Id.String(),
			Username: user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(p.secret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenStr, nil
}

func (p *jwtTokenProvider) RenewAccessToken(refreshTokenString string) (*string, error) {
	// Parse and verify the refresh token
	refreshToken, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})

	if err != nil || !refreshToken.Valid {
		return nil, errs.InvalidToken
	}

	// Generate a new access token if refresh token is valid
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		username := claims["username"].(string)
		userId := claims["user_id"].(string)

		parsedUUID, err := uuid.Parse(userId)
		if err != nil {
			return nil, errs.ParseUUIDError
		}

		newAccessTokenString, err := p.GenerateAccessToken(model.User{Username: username, Id: parsedUUID})
		if err != nil {
			return nil, err
		}

		return &newAccessTokenString, nil
	} else {
		return nil, errs.InvalidToken
	}
}

func (p *jwtTokenProvider) ExtractToken(authHeader string) (string, error) {
	splits := strings.Split(authHeader, " ")

	if len(splits) != 2 || splits[0] != "Bearer" {
		return "", errs.InvalidBearerFormat
	}

	return splits[1], nil
}

func (p *jwtTokenProvider) ValidateToken(token string) (*JwtClaims, error) {
	claims := JwtClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})

	if jwtToken == nil || !jwtToken.Valid {
		return nil, errs.InvalidToken
	}

	if err != nil {
		return nil, err
	}

	if claims.Issuer != p.issuer {
		return nil, errs.InvalidIssuer
	}

	return &claims, nil
}
