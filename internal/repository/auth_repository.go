package repository

import (
	"time"

	"github.com/EputraP/kfc_be/internal/dto"
	"github.com/EputraP/kfc_be/internal/model"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"gorm.io/gorm"
)

type AuthRepository interface {
	WithTx(tx *gorm.DB) AuthRepository
	CreateUser(input *dto.RegisterBody) (*model.User, error)
	SearchUserByUsername(input *dto.RegisterBody) (*model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
func (r authRepository) WithTx(tx *gorm.DB) AuthRepository {
	return &authRepository{
		db: tx,
	}
}

func (r *authRepository) CreateUser(input *dto.RegisterBody) (*model.User, error) {

	logger.Info("userRepository CreateUser", "Executing CreateUser SQL query", map[string]string{
		"username": input.Username,
	})

	resultModel := &model.User{}

	sqlScript := `INSERT INTO users (username, password, created_at) 
				VALUES (?,?,?) 
				RETURNING id, username;`

	res := r.db.Raw(sqlScript, input.Username, input.Password, time.Now()).Scan(resultModel)

	if res.Error != nil {
		logger.Error("userRepository CreateUser", "Failed to create user", map[string]string{
			"username": input.Username,
		})
		return nil, res.Error
	}

	logger.Info("userRepository CreateUser", "Successfully created user in CreateUser", map[string]string{
		"username": input.Username,
	})

	return resultModel, nil
}

func (r *authRepository) SearchUserByUsername(input *dto.RegisterBody) (*model.User, error) {

	logger.Info("userRepository SearchUserByUsername", "Executing SearchUserByUsername SQL query", map[string]string{
		"username": input.Username,
	})

	resultModel := &model.User{}

	sqlScript := `SELECT id, username
				  FROM
					users u 
				  WHERE
					username = ?;`

	res := r.db.Raw(sqlScript, input.Username).Scan(resultModel)

	if res.Error != nil {
		logger.Error("userRepository SearchUserByUsername", "Failed to search user", map[string]string{
			"username": input.Username,
		})
		return nil, res.Error
	}

	logger.Info("userRepository SearchUserByUsername", "Successfully ran SearchUserByUsername", map[string]string{
		"username": input.Username,
	})

	return resultModel, nil
}
