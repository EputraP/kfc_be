package main

import (
	"fmt"
	"log"
	"os"

	dbstore "github.com/EputraP/kfc_be/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("error loading env", err)
		log.Fatalln("error loading env", err)
	}

	prepare()

	srv := gin.Default()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if err := srv.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Println("Error running gin server: ", err)
		log.Fatalln("Error running gin server: ", err)
	}

}

func prepare() {
	_ = dbstore.Get()

	return

}
