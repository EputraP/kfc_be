package service

import (
	"regexp"
	"strings"

	"github.com/EputraP/kfc_be/internal/dto"
	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/model"
	"github.com/EputraP/kfc_be/internal/repository"
	"github.com/EputraP/kfc_be/internal/util/hasher"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"github.com/EputraP/kfc_be/internal/util/tokenprovider"
	"gorm.io/gorm"
)

type AuthService interface {
	CreateUser(input *dto.RegisterBody) (*dto.RegisterResponse, error)
	Login(input *dto.LoginBody) (*dto.LoginResponse, error)
}
type authService struct {
	authRepo    repository.AuthRepository
	hasher      hasher.Hasher
	jtwProvider tokenprovider.JWTTokenProvider
}

type AuthServiceConfig struct {
	AuthRepo    repository.AuthRepository
	Hasher      hasher.Hasher
	JwtProvider tokenprovider.JWTTokenProvider
}

func NewAuthService(config AuthServiceConfig) AuthService {
	return &authService{
		authRepo:    config.AuthRepo,
		hasher:      config.Hasher,
		jtwProvider: config.JwtProvider,
	}
}

func (s *authService) CreateUser(input *dto.RegisterBody) (*dto.RegisterResponse, error) {
	logger.Info("authService CreateUser", "Executing CreateUser Service", map[string]string{
		"username": input.Username,
	})

	lowerUsername := strings.ToLower(input.Username)

	userData, err := s.authRepo.SearchUserByUsername(&dto.RegisterBody{Username: lowerUsername})
	if err != nil {
		logger.Error("authService CreateUser", errs.SearchUsernameError.Error(), map[string]string{
			"userName": input.Username,
			"error":    err.Error(),
		})
		return nil, errs.SearchUsernameError
	}
	if len(userData.Username) != 0 {
		logger.Error("authService CreateUser", errs.UsernameAlreadyUsed.Error(), map[string]string{
			"userName": input.Username,
		})
		return nil, errs.UsernameAlreadyUsed
	}

	re := regexp.MustCompile(`(?i)` + input.Username)
	isMatch := re.MatchString(input.Password)

	if isMatch {
		logger.Error("authService CreateUser", errs.PasswordContainUsername.Error(), map[string]string{
			"userName": input.Username,
		})
		return nil, errs.PasswordContainUsername
	}

	resp := &dto.RegisterResponse{}

	err = repository.AsTransaction(func(tx *gorm.DB) error {
		repoWithTx := s.authRepo.WithTx(tx)

		hashedPassword, _ := s.hasher.Hash(input.Password)

		newUser, err := repoWithTx.CreateUser(&dto.RegisterBody{
			Username: lowerUsername,
			Password: hashedPassword,
		})
		if err != nil {
			logger.Error("authService CreateUser", "Error creating new user", map[string]string{
				"userName": input.Username,
				"error":    err.Error(),
			})
			return err
		}

		resp = &dto.RegisterResponse{
			UserID:   newUser.Id,
			Username: newUser.Username,
		}

		return nil
	})

	if err != nil {
		logger.Error("authService CreateUser", "Error transaction", map[string]string{
			"userName": input.Username,
			"error":    err.Error(),
		})
		return nil, err
	}

	logger.Info("authService CreateUser", "Finished CreateUser Service", map[string]string{
		"username": input.Username,
	})

	return resp, nil
}

func (s authService) Login(input *dto.LoginBody) (*dto.LoginResponse, error) {

	logger.Info("authService Login", "Executing Login Service", map[string]string{
		"username": input.Username,
	})

	lowerUsername := strings.ToLower(input.Username)

	account, err := s.authRepo.SearchUserByUsername(&dto.RegisterBody{Username: lowerUsername})
	if err != nil {
		logger.Error("authService Login", errs.SearchUsernameError.Error(), map[string]string{
			"userName": input.Username,
			"error":    err.Error(),
		})
		return nil, errs.SearchUsernameError
	}
	if len(account.Username) == 0 {
		logger.Error("authService CreateUser", errs.UsernamePasswordIncorrect.Error(), map[string]string{
			"userName": input.Username,
		})
		return nil, errs.UsernamePasswordIncorrect
	}

	passwordOk, err := s.hasher.IsEqual(account.Password, input.Password)
	if err != nil {
		logger.Error("authService CreateUser", errs.CheckPasswordError.Error(), map[string]string{
			"userName": input.Username,
			"error":    err.Error(),
		})
		return nil, errs.CheckPasswordError
	}

	if !passwordOk {
		logger.Error("authService CreateUser", errs.PasswordDoesntMatch.Error(), map[string]string{
			"userName": input.Username,
		})
		return nil, errs.PasswordDoesntMatch
	}

	loginResponse, err := s.generateLoginResponse(&model.User{Id: account.Id, Username: account.Username})
	if err != nil {
		logger.Error("authService CreateUser", errs.GenerateLoginResponseError.Error(), map[string]string{
			"userName": input.Username,
			"error":    err.Error(),
		})
		return nil, errs.GenerateLoginResponseError
	}

	logger.Info("authService Login", "Finished Login Service", map[string]string{
		"username": input.Username,
	})

	return loginResponse, nil
}

func (as authService) generateLoginResponse(user *model.User) (*dto.LoginResponse, error) {
	accesToken, err := as.jtwProvider.GenerateAccessToken(*user)

	if err != nil {
		return nil, err
	}

	refreshToken, err := as.jtwProvider.GenerateRefreshToken(*user)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccesToken:   accesToken,
		RefreshToken: refreshToken,
	}, nil
}
