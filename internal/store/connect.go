package dbstore

import (
	"fmt"
	"log"
	"os"
	"sync"

	logs "github.com/EputraP/kfc_be/internal/util/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once sync.Once
	db   *gorm.DB
)

func connectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, pass, dbName)

	log.Println("Connecting with DSN: ", dsn)
	logs.Info("connectDB", "Connecting with DSN: ", map[string]string{
		"dsn": dsn,
	})
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logs.Error("connectDB", "error gorm.Open: ", map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return dbConn, err
}

func connect() {
	dbConn, err := connectDB()

	if err != nil {
		logs.Error("connect", "error dbConn.DB: ", map[string]string{
			"error": err.Error(),
		})
	}

	logs.Info("connectDB", "Success connecting to db", nil)

	db = dbConn
}

func Get() *gorm.DB {
	once.Do(func() {
		connect()
	})

	return db
}
