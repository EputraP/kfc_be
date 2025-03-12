package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `json:"username" gorm:"type:varchar;not null"`
	Password string    `json:"password" gorm:"type:varchar;not null"`
}
