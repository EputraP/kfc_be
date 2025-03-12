package service

import (
	"regexp"
	"strings"

	"github.com/EputraP/kfc_be/internal/dto"
	errs "github.com/EputraP/kfc_be/internal/errors"
	"github.com/EputraP/kfc_be/internal/repository"
	"github.com/EputraP/kfc_be/internal/util/hasher"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"gorm.io/gorm"
)

type AuthService interface {
	CreateUser(input *dto.RegisterBody) (*dto.RegisterResponse, error)
}
type authService struct {
	authRepo repository.AuthRepository
	hasher   hasher.Hasher
}

type AuthServiceConfig struct {
	AuthRepo repository.AuthRepository
	Hasher   hasher.Hasher
}

func NewAuthService(config AuthServiceConfig) AuthService {
	return &authService{
		authRepo: config.AuthRepo,
		hasher:   config.Hasher,
	}
}

func (s *authService) CreateUser(input *dto.RegisterBody) (*dto.RegisterResponse, error) {
	lowerUsername := strings.ToLower(input.Username)

	userData, err := s.authRepo.SearchUserByUsername(&dto.RegisterBody{Username: lowerUsername})
	if err != nil {
		logger.Error("authService CreateUser", errs.EmailAlreadyUsed.Error(), map[string]string{
			"userName": input.Username,
		})
		return nil, errs.EmailAlreadyUsed
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

	return resp, nil
}
