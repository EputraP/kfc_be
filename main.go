package main

import (
	"fmt"
	"log"
	"os"

	dbstore "github.com/EputraP/kfc_be/internal/store"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func main() {

	if err := logger.Init("app.log"); err != nil {
		logger.Error("main", "Failed to initialize logger:", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Info("main", "Starting application...", nil)

	err := godotenv.Load()
	if err != nil {
		logger.Error("main", "Error loading .env file", map[string]string{
			"error": err.Error(),
		})
		return
	}

	prepare()

	srv := gin.Default()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Info("main",
		"Server is starting...", map[string]string{
			"port": port,
		})

	if err := srv.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Println("Error running gin server: ", err)
		log.Fatalln("Error running gin server: ", err)
	}

}

func prepare() {
	logger.Info("main", "Initializing dependencies...", nil)

	_ = dbstore.Get()

	logger.Info("main", "Initializing repositories...", nil)

	logger.Info("main", "Initializing services...", nil)

	logger.Info("main", "Initializing handlers...", nil)

	logger.Info("main", "Application initialized successfully.", nil)

}
