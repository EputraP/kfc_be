package repository

import (
	dbstore "github.com/EputraP/kfc_be/internal/store"
	"gorm.io/gorm"
)

type TransactionFunc func(tx *gorm.DB) error

func AsTransaction(transactionFn TransactionFunc) error {
	db := dbstore.Get()
	tx := db.Begin()

	if err := transactionFn(tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
