package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akshaypatil3096/url-shortener/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Printf("failed to load env, err: %v", err))
	}

	router := gin.Default()
	router.GET("/:url")
	router.POST("/app/v1", controller.ShortenerURL)
	if err := router.Run(os.Getenv("DOMAIN")); err != nil {
		log.Fatal(fmt.Printf("failed to start the service, err: %v", err))

	}

}
